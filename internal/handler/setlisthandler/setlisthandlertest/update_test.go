package setlisthandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/domain/mocks"
	"github.com/96Asch/mkvstage-server/internal/handler/setlisthandler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/datatypes"
)

func prepareAndServeUpdate(
	t *testing.T,
	param string,
	mockSL domain.SetlistService,
	mockSS domain.SongService,
	mockMWH domain.MiddlewareHandler,
	body *[]byte,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	writer := httptest.NewRecorder()

	setlisthandler.Initialize(&router.RouterGroup, mockSL, mockSS, mockMWH)

	requestBody := bytes.NewReader(*body)

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPut,
		fmt.Sprintf("/setlists/%s/update", param),
		requestBody,
	)
	assert.NoError(t, err)

	router.ServeHTTP(writer, req)

	return writer
}

func TestUpdateByIDCorrect(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	mockPrevSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Bar",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().Add(time.Hour * 24).Round(0),
		UpdatedAt: time.Now().Round(0),
		Order:     datatypes.JSON([]byte(`{"order" : "2,1,3,4"}`)),
	}

	expSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().Add(time.Hour * 24).Round(0),
		UpdatedAt: time.Now().Round(0),
		Order:     datatypes.JSON([]byte(`{"order" : "1,2,3,4"}`)),
	}

	mockSetlist := &domain.Setlist{
		Name:      expSetlist.Name,
		CreatorID: expSetlist.ID,
		Deadline:  expSetlist.Deadline,
	}

	expMockSetlist := &domain.Setlist{
		ID:        mockPrevSetlist.ID,
		Name:      expSetlist.Name,
		CreatorID: expSetlist.ID,
		Deadline:  expSetlist.Deadline,
	}

	mockSL := &mocks.MockSetlistService{}
	mockSS := &mocks.MockSongService{}

	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	mockSL.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), expMockSetlist, mockUser).
		Return(expSetlist, nil)

	byteBody, err := json.Marshal(mockSetlist)
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, fmt.Sprint(mockPrevSetlist.ID), mockSL, mockSS, mockMWH, &byteBody)

	assert.Equal(t, http.StatusOK, writer.Code)

	expBody, err := json.Marshal(gin.H{"setlist": expSetlist})
	assert.NoError(t, err)

	assert.Equal(t, expBody, writer.Body.Bytes())
	mockSL.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDInvalidParam(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:         1,
		Permission: domain.GUEST,
	}

	expSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().Add(time.Hour * 24).Round(0),
		UpdatedAt: time.Now().Round(0),
		Order:     datatypes.JSON([]byte(`{"order" : "1,2,3,4"}`)),
	}

	mockSetlist := &domain.Setlist{
		Name:      expSetlist.Name,
		CreatorID: expSetlist.ID,
		Deadline:  expSetlist.Deadline,
	}

	mockMWH := &mocks.MockMiddlewareHandler{}
	mockSL := &mocks.MockSetlistService{}
	mockSS := &mocks.MockSongService{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(mockSetlist)
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, "A", mockSL, mockSS, mockMWH, &byteBody)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
	mockMWH.AssertExpectations(t)
	mockSL.AssertExpectations(t)
	mockSS.AssertExpectations(t)
}

func TestUpdateByIDBindErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	mockPrevSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Bar",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().Add(time.Hour * 24).Round(0),
		UpdatedAt: time.Now().Round(0),
		Order:     datatypes.JSON([]byte(`{"order" : "2,1,3,4"}`)),
	}

	mockSetlist := &domain.Setlist{
		CreatorID: mockUser.ID,
	}

	mockErr := domain.NewBadRequestErr("")
	mockSL := &mocks.MockSetlistService{}
	mockSS := &mocks.MockSongService{}

	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(mockSetlist)
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, fmt.Sprint(mockPrevSetlist.ID), mockSL, mockSS, mockMWH, &byteBody)

	assert.Equal(t, mockErr.Status(), writer.Code)
	mockSL.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDNoContext(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	mockPrevSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Bar",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().Add(time.Hour * 24).Round(0),
		UpdatedAt: time.Now().Round(0),
		Order:     datatypes.JSON([]byte(`{"order" : "2,1,3,4"}`)),
	}

	mockSetlist := &domain.Setlist{
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().Add(time.Hour * 24),
	}

	mockErr := domain.NewInternalErr()
	mockSL := &mocks.MockSetlistService{}
	mockSS := &mocks.MockSongService{}

	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	byteBody, err := json.Marshal(mockSetlist)
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, fmt.Sprint(mockPrevSetlist.ID), mockSL, mockSS, mockMWH, &byteBody)

	assert.Equal(t, mockErr.Status(), writer.Code)
	mockSL.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}

func TestUpdateByIDStoreErr(t *testing.T) {
	t.Parallel()

	mockUser := &domain.User{
		ID:        1,
		FirstName: "Foo",
		LastName:  "Bar",
	}

	mockPrevSetlist := &domain.Setlist{
		ID:        1,
		Name:      "Bar",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().Add(time.Hour * 24).Round(0),
		UpdatedAt: time.Now().Round(0),
		Order:     datatypes.JSON([]byte(`{"order" : "2,1,3,4"}`)),
	}

	mockSetlist := &domain.Setlist{
		Name:      "Foo",
		CreatorID: mockUser.ID,
		Deadline:  time.Now().Add(time.Hour * 24).Round(0),
	}

	expMockSetlist := mockSetlist
	expMockSetlist.ID = mockPrevSetlist.ID

	mockErr := domain.NewBadRequestErr("")
	mockSL := &mocks.MockSetlistService{}
	mockSS := &mocks.MockSongService{}

	mockMWH := &mocks.MockMiddlewareHandler{}

	var mockAuthHF gin.HandlerFunc = func(ctx *gin.Context) {
		ctx.Set("user", mockUser)
		ctx.Next()
	}

	mockMWH.
		On("AuthenticateUser").
		Return(mockAuthHF)

	mockSL.
		On("Update", mock.AnythingOfType("*context.emptyCtx"), expMockSetlist, mockUser).
		Return(nil, mockErr)

	byteBody, err := json.Marshal(mockSetlist)
	assert.NoError(t, err)

	writer := prepareAndServeUpdate(t, fmt.Sprint(mockPrevSetlist.ID), mockSL, mockSS, mockMWH, &byteBody)

	assert.Equal(t, mockErr.Status(), writer.Code)
	mockSL.AssertExpectations(t)
	mockMWH.AssertExpectations(t)
}
