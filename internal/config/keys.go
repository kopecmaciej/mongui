package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/kopecmaciej/vi-mongo/internal/util"
)

type (
	OrderedKeys struct {
		Element string
		Keys    []Key
	}
	// KeyBindings is a way to define keybindings for the application
	// There are views that have only keybindings and some have
	// nested keybindings of their children views
	KeyBindings struct {
		Global    GlobalKeys    `json:"global"`
		Help      HelpKeys      `json:"help"`
		Welcome   WelcomeKeys   `json:"welcome"`
		Connector ConnectorKeys `json:"connector"`
		Root      RootKeys      `json:"root"`
		Databases DatabasesKeys `json:"databases"`
		Content   ContentKeys   `json:"content"`
		DocPeeker DocPeekerKeys `json:"docPeeker"`
		History   HistoryKeys   `json:"history"`
	}

	// Key is a lowest level of keybindings
	// It holds the keys and runes that are used to trigger the action
	// and a description of the action that will be displayed in the help
	Key struct {
		Keys        []string `json:"keys,omitempty"`
		Runes       []string `json:"runes,omitempty"`
		Description string   `json:"description"`
	}

	// GlobalKeys is a struct that holds the global keybindings
	// for the application, they can be triggered from any view
	// as keys are passed from top to bottom
	GlobalKeys struct {
		ToggleFullScreenHelp Key `json:"toggleFullScreenHelp"`
	}

	RootKeys struct {
		ToggleFocus    Key `json:"toggleFocus"`
		FocusDatabases Key `json:"focusDatabases"`
		FocusContent   Key `json:"focusContent"`
		HideDatabases  Key `json:"hideDatabases"`
		OpenConnector  Key `json:"openConnector"`
	}

	DatabasesKeys struct {
		FilterBar        Key `json:"filterBar"`
		ExpandAll        Key `json:"expandAll"`
		CollapseAll      Key `json:"collapseAll"`
		ToggleExpand     Key `json:"toggleExpand"`
		AddCollection    Key `json:"addCollection"`
		DeleteCollection Key `json:"deleteCollection"`
	}

	ContentKeys struct {
		ChangeView        Key      `json:"switchView"`
		PeekDocument      Key      `json:"peekDocument"`
		ViewDocument      Key      `json:"viewDocument"`
		AddDocument       Key      `json:"addDocument"`
		EditDocument      Key      `json:"editDocument"`
		DuplicateDocument Key      `json:"duplicateDocument"`
		DeleteDocument    Key      `json:"deleteDocument"`
		MultipleSelect    Key      `json:"multipleSelect"`
		ClearSelection    Key      `json:"clearSelection"`
		CopyLine          Key      `json:"copyValue"`
		Refresh           Key      `json:"refresh"`
		ToggleQuery       Key      `json:"toggleQuery"`
		NextDocument      Key      `json:"nextDocument"`
		PreviousDocument  Key      `json:"previousDocument"`
		NextPage          Key      `json:"nextPage"`
		PreviousPage      Key      `json:"previousPage"`
		QueryBar          QueryBar `json:"queryBar"`
		ToggleSort        Key      `json:"toggleSort"`
	}

	QueryBar struct {
		ShowHistory Key `json:"showHistory"`
		ClearInput  Key `json:"clearInput"`
	}

	ConnectorKeys struct {
		ToggleFocus   Key               `json:"toggleFocus"`
		ConnectorForm ConnectorFormKeys `json:"connectorForm"`
		ConnectorList ConnectorListKeys `json:"connectorList"`
	}

	ConnectorFormKeys struct {
		SaveConnection Key `json:"saveConnection"`
		FocusList      Key `json:"focusList"`
	}

	ConnectorListKeys struct {
		FocusForm        Key `json:"focusForm"`
		DeleteConnection Key `json:"deleteConnection"`
		SetConnection    Key `json:"setConnection"`
	}

	WelcomeKeys struct {
		MoveFocusUp   Key `json:"moveFocusUp"`
		MoveFocusDown Key `json:"moveFocusDown"`
	}

	HelpKeys struct {
		Close Key `json:"close"`
	}

	DocPeekerKeys struct {
		MoveToTop    Key `json:"moveToTop"`
		MoveToBottom Key `json:"moveToBottom"`
		CopyFullObj  Key `json:"copyFullObj"`
		CopyValue    Key `json:"copyValue"`
		Refresh      Key `json:"refresh"`
	}

	HistoryKeys struct {
		ClearHistory Key `json:"clearHistory"`
		AcceptEntry  Key `json:"acceptEntry"`
		CloseHistory Key `json:"closeHistory"`
	}
)

func (k *KeyBindings) loadDefaultKeybindings() {
	k.Global = GlobalKeys{
		ToggleFullScreenHelp: Key{
			Runes:       []string{"?"},
			Description: "Toggle full screen help",
		},
	}

	k.Root = RootKeys{
		ToggleFocus: Key{
			Keys:        []string{"Tab", "Backtab"},
			Description: "Focus next view",
		},
		FocusDatabases: Key{
			Keys:        []string{"Ctrl+H"},
			Description: "Focus databases",
		},
		FocusContent: Key{
			Keys:        []string{"Ctrl+L"},
			Description: "Focus content",
		},
		HideDatabases: Key{
			Keys:        []string{"Ctrl+N"},
			Description: "Hide databases",
		},
		OpenConnector: Key{
			Keys:        []string{"Ctrl+O"},
			Description: "Open connector",
		},
	}

	k.Databases = DatabasesKeys{
		FilterBar: Key{
			Runes:       []string{"/"},
			Description: "Focus filter bar",
		},
		ExpandAll: Key{
			Runes:       []string{"E"},
			Description: "Expand all",
		},
		CollapseAll: Key{
			Runes:       []string{"W"},
			Description: "Collapse all",
		},
		ToggleExpand: Key{
			Runes:       []string{"T"},
			Description: "Toggle expand",
		},
		AddCollection: Key{
			Runes:       []string{"A"},
			Description: "Add collection",
		},
		DeleteCollection: Key{
			Runes:       []string{"D"},
			Description: "Delete collection",
		},
	}

	k.Content = ContentKeys{
		ChangeView: Key{
			Runes:       []string{"f"},
			Description: "Change table view",
		},
		PeekDocument: Key{
			Runes:       []string{"p"},
			Keys:        []string{"Enter"},
			Description: "Peek document",
		},
		ViewDocument: Key{
			Runes:       []string{"P"},
			Description: "View document",
		},
		AddDocument: Key{
			Runes:       []string{"A"},
			Description: "Add document",
		},
		EditDocument: Key{
			Runes:       []string{"E"},
			Description: "Edit document",
		},
		DuplicateDocument: Key{
			Runes:       []string{"d"},
			Description: "Duplicate document",
		},
		DeleteDocument: Key{
			Runes:       []string{"D"},
			Description: "Delete document",
		},
		MultipleSelect: Key{
			Runes:       []string{"v"},
			Description: "Multiple select",
		},
		ClearSelection: Key{
			Runes:       []string{"C"},
			Description: "Clear selection",
		},
		CopyLine: Key{
			Runes:       []string{"c"},
			Description: "Copy value",
		},
		Refresh: Key{
			Runes:       []string{"R"},
			Description: "Refresh",
		},
		ToggleQuery: Key{
			Runes:       []string{"/"},
			Description: "Toggle query",
		},
		ToggleSort: Key{
			Runes:       []string{"s"},
			Description: "Toggle sort",
		},
		NextDocument: Key{
			Runes:       []string{"]"},
			Description: "Next document",
		},
		PreviousDocument: Key{
			Runes:       []string{"["},
			Description: "Previous document",
		},
		NextPage: Key{
			Runes:       []string{"n"},
			Description: "Next page",
		},
		PreviousPage: Key{
			Runes:       []string{"b"},
			Description: "Previous page",
		},
	}

	k.Content.QueryBar = QueryBar{
		ShowHistory: Key{
			Keys:        []string{"Ctrl+Y"},
			Description: "Show history",
		},
		ClearInput: Key{
			Keys:        []string{"Ctrl+D"},
			Description: "Clear input",
		},
	}

	k.Connector.ToggleFocus = Key{
		Keys:        []string{"Tab", "Backtab"},
		Description: "Toggle focus",
	}

	k.Connector.ConnectorForm = ConnectorFormKeys{
		SaveConnection: Key{
			Keys:        []string{"Ctrl+S", "Enter"},
			Description: "Save connection",
		},
		FocusList: Key{
			Keys:        []string{"Ctrl+H", "Ctrl+Left"},
			Description: "Focus Connection List",
		},
	}

	k.Connector.ConnectorList = ConnectorListKeys{
		FocusForm: Key{
			Keys:        []string{"Ctrl+L", "Ctrl+Right"},
			Description: "Move focus to form",
		},
		DeleteConnection: Key{
			Runes:       []string{"D"},
			Description: "Delete selected connection",
		},
		SetConnection: Key{
			Keys:        []string{"Enter", "Space"},
			Description: "Set selected connection",
		},
	}

	k.Welcome = WelcomeKeys{
		MoveFocusUp: Key{
			Keys:        []string{"Backtab"},
			Description: "Move focus up",
		},
		MoveFocusDown: Key{
			Keys:        []string{"Tab"},
			Description: "Move focus down",
		},
	}

	k.Help = HelpKeys{
		Close: Key{
			Keys:        []string{"Esc"},
			Description: "Close help",
		},
	}

	k.DocPeeker = DocPeekerKeys{
		MoveToTop: Key{
			Runes:       []string{"g"},
			Description: "Move to top",
		},
		MoveToBottom: Key{
			Runes:       []string{"G"},
			Description: "Move to bottom",
		},
		CopyFullObj: Key{
			Runes:       []string{"c"},
			Description: "Copy full object",
		},
		CopyValue: Key{
			Runes:       []string{"v"},
			Description: "Copy value",
		},
		Refresh: Key{
			Runes:       []string{"R"},
			Description: "Refresh document",
		},
	}

	k.History = HistoryKeys{
		ClearHistory: Key{
			Runes:       []string{"C"},
			Description: "Clear history",
		},
		AcceptEntry: Key{
			Keys:        []string{"Enter", "Space"},
			Description: "Accept entry",
		},
		CloseHistory: Key{
			Keys:        []string{"Esc", "Ctrl+Y"},
			Description: "Close history",
		},
	}
}

// LoadKeybindings loads keybindings from the config file
// if the file does not exist it creates a new one with default keybindings
func LoadKeybindings() (*KeyBindings, error) {
	keybindings := &KeyBindings{}
	defaultKeybindings := &KeyBindings{}
	defaultKeybindings.loadDefaultKeybindings()

	keybindingsPath, err := getKeybindingsPath()
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(keybindingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			// just for easy development
			if os.Getenv("ENV") == "vi-dev" {
				return defaultKeybindings, nil
			}
			err := ensureConfigDirExist()
			if err != nil {
				return nil, err
			}
			bytes, err = json.Marshal(defaultKeybindings)
			if err != nil {
				return nil, err
			}
			err = os.WriteFile(keybindingsPath, bytes, 0644)
			if err != nil {
				return nil, err
			}
			return defaultKeybindings, nil
		}
		return nil, err
	}

	err = json.Unmarshal(bytes, keybindings)
	if err != nil {
		return nil, err
	}

	util.MergeConfigs(keybindings, defaultKeybindings)

	return keybindings, nil
}

// extractKeysFromStruct extracts all Key structs from a reflect.Value
func extractKeysFromStruct(val reflect.Value) []Key {
	var keys []Key

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Type() == reflect.TypeOf(Key{}) {
			keys = append(keys, field.Interface().(Key))
		} else if field.Kind() == reflect.Struct {
			keys = append(keys, extractKeysFromStruct(field)...)
		}
	}

	return keys
}

// GetAvaliableKeys returns all keys
func (kb KeyBindings) GetAvaliableKeys() []OrderedKeys {
	var keys []OrderedKeys

	v := reflect.ValueOf(kb)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name

		orderedKeys := OrderedKeys{
			Element: fieldName,
			Keys:    extractKeysFromStruct(field),
		}

		keys = append(keys, orderedKeys)
	}

	return keys
}

// GetKeysForElement returns keys for element
func (kb KeyBindings) GetKeysForElement(elementId string) ([]OrderedKeys, error) {
	if elementId == "" {
		return nil, fmt.Errorf("element is empty")
	}

	v := reflect.ValueOf(kb)
	field := v.FieldByName(elementId)

	if !field.IsValid() || field.Kind() != reflect.Struct {
		return nil, fmt.Errorf("field %s not found", elementId)
	}

	keys := []OrderedKeys{{
		Element: elementId,
		Keys:    extractKeysFromStruct(field),
	}}

	return keys, nil
}

// ConvertStrKeyToTcellKey converts string key to tcell key
func (kb *KeyBindings) ConvertStrKeyToTcellKey(key string) (tcell.Key, bool) {
	for k, v := range tcell.KeyNames {
		if v == key {
			return k, true
		}
	}
	return -1, false
}

// Contains checks if the keybindings contains the key
func (kb *KeyBindings) Contains(configKey Key, namedKey string) bool {
	// some hacks for couple of keys
	if namedKey == "Rune[ ]" {
		namedKey = "Space"
	}
	// in some terminals ctrl+H often is seen as backspace
	if namedKey == "Backspace" {
		namedKey = "Ctrl+H"
	}

	if strings.HasPrefix(namedKey, "Rune") {
		namedKey = strings.TrimPrefix(namedKey, "Rune")
		for _, k := range configKey.Runes {
			if k == namedKey[1:2] {
				return true
			}
		}
	}

	for _, k := range configKey.Keys {
		if k == namedKey {
			return true
		}
	}

	return false
}

func (k *Key) String() string {
	var keyString string
	var iter []string
	if len(k.Keys) > 0 {
		iter = k.Keys
	} else {
		iter = k.Runes
	}
	for i, k := range iter {
		if i == 0 {
			keyString = k
		} else {
			keyString = fmt.Sprintf("%s, %s", keyString, k)
		}
	}

	return keyString
}

func getKeybindingsPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return configDir + "/keybindings.json", nil
}
