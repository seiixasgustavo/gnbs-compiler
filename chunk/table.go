package chunk

const TableMaxLoad = 0.75

type Table struct {
	Entries []*Entry
	Count   int
}

type Entry struct {
	Key   *GString
	Value Value
}

func NewTable() *Table {
	return &Table{nil, 0}
}

func FindEntry(entries []*Entry, key *GString) *Entry {
	capacity := uint32(cap(entries))
	index := key.Hash % capacity

	var tombstone *Entry = nil

	for {
		if entry := entries[index]; entry.Key == nil {
			if entry.Value.Type == TypeNull {
				if tombstone == nil {
					return entry
				} else {
					return tombstone
				}
			} else {
				if tombstone == nil {
					tombstone = entry
				}
			}
		} else if entry.Key == key {
			return entry
		}

		index = (index + 1) % capacity
	}
}

func (t *Table) AdjustCapacity() {
	capacity := cap(t.Entries)
	if capacity < 8 {
		capacity = 8
	} else {
		capacity *= 2
	}

	entries := make([]*Entry, capacity)

	for i := 0; i < capacity; i++ {
		entries[i].Key = nil
		entries[i].Value = Value{}
	}

	t.Count = 0
	for i := 0; i < cap(t.Entries); i++ {
		entry := t.Entries[i]
		if entry == nil {
			continue
		}
		dest := FindEntry(entries, entry.Key)
		dest.Key = entry.Key
		dest.Value = entry.Value
		t.Count++
	}
	t.Entries = entries
}

func (t *Table) TableSet(key *GString, value Value) bool {
	if float64(t.Count+1) > float64(cap(t.Entries))*TableMaxLoad {
		t.AdjustCapacity()
	}

	entry := FindEntry(t.Entries, key)

	isNewKey := entry.Key == nil
	if isNewKey && entry.Value.Type == TypeNull {
		t.Count++
	}
	entry.Key = key
	entry.Value = value
	return isNewKey
}

func (t *Table) TableGet(key *GString, value *Value) bool {
	if t.Count == 0 {
		return false
	}
	entry := FindEntry(t.Entries, key)
	if entry.Key == nil {
		return false
	}
	*value = entry.Value
	return true
}

func (t *Table) TableDelete(key *GString) bool {
	if t.Count == 0 {
		return false
	}

	entry := FindEntry(t.Entries, key)
	if entry.Key == nil {
		return false
	}

	entry.Key = nil
	entry.Value = Value{
		Type:  TypeBool,
		Value: false,
	}

	return true
}

func TableAddAll(from, to *Table) {
	for i := 0; i < cap(from.Entries); i++ {
		if entry := from.Entries[i]; entry.Key != nil {
			to.TableSet(entry.Key, entry.Value)
		}
	}
}
func TableFindString(table *Table, message string, hash uint32) *GString {
	if table.Count == 0 {
		return nil
	}

	index := hash % uint32(cap(table.Entries))

	for {
		entry := table.Entries[index]
		if entry.Key == nil {
			if entry.Value.Type == TypeNull {
				return nil
			}
		} else if len(entry.Key.String) == len(message) &&
			entry.Key.Hash == hash && entry.Key.String == message {
			return entry.Key
		}
		index = (index + 1) % uint32(cap(table.Entries))
	}
}
