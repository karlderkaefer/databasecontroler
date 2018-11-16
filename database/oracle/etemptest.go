package oracle

//
// import (
// 	"errors"
// 	"testing"
//
// 	database "github.com/karlderkaefer/databasecontroler/database"
// 	mocks "github.com/karlderkaefer/databasecontroler/database/mocks"
// 	. "github.com/stretchr/testify/mock"
// )
//
// func TestCreateUser(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	oracle := NewOracleHandler()
// 	resp, err := oracle.CreateUser("testusercreate", "testpass")
// 	t.Logf("%v", resp)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	expectMessage := "user created testusercreate"
// 	if resp.Content != expectMessage {
// 		t.Errorf("expected message: %s but was %s", expectMessage, resp.Content)
// 	}
// 	oracle.DropUser("testusercreate")
// }
//
// func TestMockDropUser(t *testing.T) {
// 	handler := mocks.DatabaseHandler{}
// 	handler.On("Execute", Anything).Return(database.Message{}, nil)
// 	oracle := Oracle{
// 		handler: handler,
// 	}
// 	msg, err := oracle.DropUser("peter")
// 	expectMessage := "user dropped peter"
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if msg.Content != expectMessage {
// 		t.Errorf("expected message: %s but was %s", expectMessage, msg.Content)
// 	}
//
// 	handlerConFailed := mocks.DatabaseHandler{}
// 	handlerConFailed.On("Execute", Anything).Return(database.Message{}, errors.New("No Connection"))
// 	oracle = Oracle{
// 		handler: handlerConFailed,
// 	}
// 	msg, err = oracle.DropUser("hello")
// 	expectMessage = "No Connection"
// 	if err.Error() != expectMessage {
// 		t.Errorf("expected message: %s but was %s", expectMessage, msg.Content)
// 	}
// }
//
// func TestMockCreateUser(t *testing.T) {
// 	handler := mocks.DatabaseHandler{}
// 	handler.On("Execute", Anything).Return(database.Message{}, nil)
// 	oracle := Oracle{
// 		handler: handler,
// 	}
// 	msg, err := oracle.CreateUser("peter", "pass")
// 	expectMessage := "user created peter"
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if msg.Content != expectMessage {
// 		t.Errorf("expected message: %s but was %s", expectMessage, msg.Content)
// 	}
// }
//
// func TestMockCreateUserAlreadyExist(t *testing.T) {
// 	handler := mocks.DatabaseHandler{}
// 	handler.On("Execute", Anything).Return(database.Message{}, errors.New("ORA-01920"))
// 	oracle := Oracle{
// 		handler: handler,
// 	}
// 	msg, err := oracle.CreateUser("peter", "pass")
// 	expectMessage := "ORA-01920"
// 	if err.Error() != expectMessage {
// 		t.Errorf("expected error: %s but was %s", expectMessage, err.Error())
// 	}
// 	if msg.Content != expectMessage {
// 		t.Errorf("expected message: %s but was %s", expectMessage, msg.Content)
// 	}
// 	if msg.Severity != database.Warn {
// 		t.Errorf("expected log level: %s but was %s", database.Warn, msg.Severity)
// 	}
// }
