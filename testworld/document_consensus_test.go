// +build testworld

package testworld

import (
	"net/http"
	"testing"
)

func TestHost_AddExternalCollaborator(t *testing.T) {

	alice := doctorFord.getHostTestSuite(t, "Alice")
	bob := doctorFord.getHostTestSuite(t, "Bob")
	charlie := doctorFord.getHostTestSuite(t, "Charlie")

	// Alice shares invoice document with Bob first
	res, err := alice.host.createInvoice(alice.expect, http.StatusOK, defaultInvoicePayload([]string{bob.id.String()}))
	if err != nil {
		t.Error(err)
	}

	docIdentifier := getDocumentIdentifier(t, res)

	if docIdentifier == "" {
		t.Error("docIdentifier empty")
	}
	params := map[string]interface{}{
		"document_id": docIdentifier,
		"currency":    "USD",
	}
	getInvoiceAndCheck(alice.expect, params)
	getInvoiceAndCheck(bob.expect, params)

	// Bob updates invoice and shares with Charlie as well
	res, err = bob.host.updateInvoice(bob.expect, http.StatusOK, docIdentifier, updatedInvoicePayload([]string{alice.id.String(), charlie.id.String()}))

	if err != nil {
		t.Error(err)
	}
	docIdentifier = getDocumentIdentifier(t, res)

	if docIdentifier == "" {
		t.Error("docIdentifier empty")
	}
	params["currency"] = "EUR"
	getInvoiceAndCheck(alice.expect, params)
	getInvoiceAndCheck(bob.expect, params)
	getInvoiceAndCheck(charlie.expect, params)
}

func TestHost_CollaboratorTimeOut(t *testing.T) {

	kenny := doctorFord.getHostTestSuite(t, "Kenny")
	bob := doctorFord.getHostTestSuite(t, "Bob")

	// Kenny shares an invoice with Bob
	response, err := kenny.host.createInvoice(kenny.expect, http.StatusOK, defaultInvoicePayload([]string{bob.id.String()}))

	if err != nil {
		t.Error(err)
	}

	// check if Bob and Kenny received the document
	docIdentifier := getDocumentIdentifier(t, response)
	paramsV1 := map[string]interface{}{
		"document_id": docIdentifier,
		"currency":    "USD",
	}
	getInvoiceAndCheck(kenny.expect, paramsV1)
	getInvoiceAndCheck(bob.expect, paramsV1)

	// Kenny gets killed
	kenny.host.canc()

	// Bob updates and sends to Alice
	updatedPayload := updatedInvoicePayload([]string{kenny.id.String()})

	// Bob will anchor the document without Alice signature but will receive an error because kenny is dead
	response, err = bob.host.updateInvoice(bob.expect, http.StatusInternalServerError, docIdentifier, updatedPayload)
	if err != nil {
		t.Error(err)
	}

	// check if bob saved the updated document
	paramsV2 := map[string]interface{}{
		"document_id": docIdentifier,
		"currency":    "EUR",
	}
	getInvoiceAndCheck(bob.expect, paramsV2)

	// bring Kenny back to life
	doctorFord.restartHost(kenny.name)

	// Kenny should not have latest version
	getInvoiceAndCheck(kenny.expect, paramsV1)

}