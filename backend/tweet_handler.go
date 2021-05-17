package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TweetInsertReq struct {
	Content string `json:"content"`
}

func (ah *AppHandlers) TweetInsertHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req TweetInsertReq
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	ids, err := ah.TweetStore.InsertChain(ctx, req.Content)
	if err != nil {
		fmt.Printf("failed insertChain: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	fmt.Fprintf(w, "done:%+v", ids)
}

type TweetUpdateReq struct {
	IDs     []string `json:"ids"`
	Content string   `json:"content"`
}

func (ah *AppHandlers) TweetUpdateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req TweetUpdateReq
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	if err := ah.TweetStore.UpdateChain(ctx, req.IDs, req.Content); err != nil {
		fmt.Printf("failed insertChain: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error:%s", err)
		return
	}
}

type TweetUpdateDMLReq struct {
	IDs     []string `json:"ids"`
	Content string   `json:"content"`
}

func (ah *AppHandlers) TweetUpdateDMLHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req TweetUpdateDMLReq
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	ids, err := ah.TweetStore.UpdateDMLChain(ctx, req.IDs, req.Content)
	if err != nil {
		fmt.Printf("failed updateDMLChain: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	fmt.Fprintf(w, "done:%+v", ids)
}

type TweetUpdateBatchDMLReq struct {
	IDs     []string `json:"ids"`
	Content string   `json:"content"`
}

func (ah *AppHandlers) TweetUpdateBatchDMLHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req TweetUpdateBatchDMLReq
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	ids, err := ah.TweetStore.UpdateBatchDMLChain(ctx, req.IDs, req.Content)
	if err != nil {
		fmt.Printf("failed updateDMLChain: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(ids); err != nil {
		fmt.Printf("failed json.Encode: %s\n", err)
	}
}

type TweetUpdateAndSelectReq struct {
	ID string `json:"id"`
}

func (ah *AppHandlers) TweetUpdateAndSelectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req TweetUpdateAndSelectReq
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	count, err := ah.TweetStore.UpdateAndSelect(ctx, req.ID)
	if err != nil {
		fmt.Printf("failed updateAndSelect: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"count": count}); err != nil {
		fmt.Printf("failed json.Encode: %s\n", err)
	}
}

type TweetUpdateDMLAndSelectReq struct {
	ID string `json:"id"`
}

func (ah *AppHandlers) TweetUpdateDMLAndSelectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req TweetUpdateDMLAndSelectReq
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	count, err := ah.TweetStore.UpdateDMLAndSelect(ctx, req.ID)
	if err != nil {
		fmt.Printf("failed updateDMLAndSelect: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"count": count}); err != nil {
		fmt.Printf("failed json.Encode: %s\n", err)
	}
}
