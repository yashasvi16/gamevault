package handler

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"strings"
	"errors"
	"time"
	"github.com/yashasvi16/gamevault/internal/model"
)

type FakePlayerRepo struct {
	ShouldFail bool
}

func (f *FakePlayerRepo) CreatePlayer(player *model.Player) error {
	if f.ShouldFail {
		return errors.New("database exploded")
	}
	player.ID = 42
	player.CreatedAt = time.Now()
	return nil
}

func (f *FakePlayerRepo) GetPlayerByEmail(email string) (*model.Player, error) {
	return nil, errors.New("not implemented")
}

func (f *FakePlayerRepo) GetLeaderboard(limit, offset int) ([]model.Player, error) {
	return nil, errors.New("not implemented")
}

func TestRegisterPlayer_Success(t *testing.T) {
	fakeRepo := &FakePlayerRepo{ShouldFail: false}
	handler := NewPlayerHandler(fakeRepo)

	body := `{"username":"tester", "email":"test@test.com","password":"longpassword123"}`
	req := httptest.NewRequest("POST", "/players", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	
	rec := httptest.NewRecorder()

	handler.RegisterPlayer(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusCreated)
	}
}

func TestRegisterPlayer_DBFailure(t *testing.T){
	fakeRepo := &FakePlayerRepo{ShouldFail: true}
	handler := NewPlayerHandler(fakeRepo)

	body := `{"username":"tester", "email":"test@test.com","password":"longpassword123"}`
	req := httptest.NewRequest("POST", "/players", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	
	rec := httptest.NewRecorder()

	handler.RegisterPlayer(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}