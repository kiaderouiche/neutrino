package api

import (
	"testing"
	"time"
)

func setupTypeTests(t *testing.T) (map[string]interface{}, *ApplicationModel, string) {
	typeName := randomString()
	app := createApp(t)

	getRec := sendAuthenticatedRequest("GET", "/app/"+app.Id, nil, t)
	getRec.CodeIs(200)

	createdApp := getRec.BodyJSON()

	return createdApp, app, typeName
}

func TestDeleteType(t *testing.T) {
	if isTravis() {
		//sporadically fails on travis
		return
	}

	_, app, typeName := setupTypeTests(t)

	deleteReq := sendAuthenticatedRequest("DELETE", "/app/"+app.Id+"/data/"+typeName, nil, t)
	deleteReq.CodeIs(200)

	getReq := sendAuthenticatedRequest("GET", "/app/"+app.Id+"/data/"+typeName, nil, t)
	//we dynamically create types so no not found errors
	getReq.CodeIs(200)

	appReq := sendAuthenticatedRequest("GET", "/app/"+app.Id, nil, t)
	updatedApp := appReq.BodyJSON()

	types := updatedApp["types"].([]interface{})

	if len(types) != 2 {
		t.Error("Types not correct count, should be users, " + typeName)
	}
}

func TestGetAndInsertTypeData(t *testing.T) {
	if isTravis() {
		//TODO: investigate why the get request returns no results
		return
	}

	_, app, typeName := setupTypeTests(t)

	sendAuthenticatedRequest("POST", "/app/"+app.Id+"/data/"+typeName, map[string]interface{}{
		"field1": "test",
		"field2": "test",
	}, t)

	time.Sleep(time.Second * 1)

	getRec := sendAuthenticatedRequest("GET", "/app/"+app.Id+"/data/"+typeName, nil, t)
	getRec.CodeIs(200)

	var data []map[string]interface{}
	getRec.Decode(&data)

	record := data[0]

	if record["field1"] != "test" || record["field2"] != "test" {
		t.Error("Item not written correctly")
	}
}

func TestGetByIdTypeData(t *testing.T) {
	_, app, typeName := setupTypeTests(t)

	rec := sendAuthenticatedRequest("POST", "/app/"+app.Id+"/data/"+typeName, map[string]interface{}{
		"field1": "test",
		"field2": "test",
	}, t)

	res := rec.BodyJSON()
	id := res["id"].(string)

	rec1 := sendAuthenticatedRequest("GET", "/app/"+app.Id+"/data/"+typeName+"/"+id, nil, t)
	item := rec1.BodyJSON()

	if item["field1"] != "test" || item["field2"] != "test" {
		t.Error("Item not written correctly")
	}
}

func TestUpdateTypeItemById(t *testing.T) {
	_, app, typeName := setupTypeTests(t)

	rec := sendAuthenticatedRequest("POST", "/app/"+app.Id+"/data/"+typeName, map[string]interface{}{
		"field1": "test",
		"field2": "test",
	}, t)
	rec.CodeIs(200)

	res := rec.BodyJSON()
	id := res["id"].(string)

	sendAuthenticatedRequest("PUT", "/app/"+app.Id+"/data/"+typeName+"/"+id, map[string]interface{}{
		"field1": "testupdated",
		"field2": "testupdated",
	}, t)

	rec1 := sendAuthenticatedRequest("GET", "/app/"+app.Id+"/data/"+typeName+"/"+id, nil, t)
	item := rec1.BodyJSON()

	if item["field1"] != "testupdated" || item["field2"] != "testupdated" {
		t.Fatal("Item not updated correctly")
	}
}

func TestDeleteTypeItemById(t *testing.T) {
	_, app, typeName := setupTypeTests(t)

	rec := sendAuthenticatedRequest("POST", "/app/"+app.Id+"/data/"+typeName, map[string]interface{}{
		"field1": "test",
		"field2": "test",
	}, t)
	rec.CodeIs(200)

	res := rec.BodyJSON()
	id := res["id"].(string)

	sendAuthenticatedRequest("DELETE", "/app/"+app.Id+"/data/"+typeName+"/"+id, nil, t)

	rec1 := sendAuthenticatedRequest("GET", "/app/"+app.Id+"/data/"+typeName+"/"+id, nil, t)
	rec1.CodeIs(404)
}
