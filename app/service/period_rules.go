package service

import (
	"fmt"
	"strings"
	"time"

	"sipasdos/app/model"
)

// File ini berisi business rules MURNI untuk periode rekrutmen:
// tidak menyentuh fiber.Ctx, tidak menyentuh database, tidak tahu apa pun
// tentang HTTP. Semua fungsinya bisa diuji tanpa menyalakan server.

const (
	dateLayout          = "2006-01-02"
	maxPeriodNameLength = 100
)

var validPeriodStatus = map[string]bool{
	model.PeriodDraft:  true,
	model.PeriodOpen:   true,
	model.PeriodClosed: true,
}

// periodTransitions menyimpan perpindahan status yang diizinkan.
var periodTransitions = map[string][]string{
	model.PeriodDraft: {model.PeriodOpen},
	model.PeriodOpen:  {model.PeriodClosed},
}

// CanTransitionPeriod memeriksa apakah status boleh berpindah dari -> ke.
// Tetap di status yang sama selalu diperbolehkan.
func CanTransitionPeriod(from, to string) bool {
	if from == to {
		return true
	}
	for _, next := range periodTransitions[from] {
		if next == to {
			return true
		}
	}
	return false
}

// PeriodSortable adalah kolom yang boleh dipakai pada ?sort=.
var PeriodSortable = map[string]bool{
	"id": true, "name": true, "start_date": true, "end_date": true, "created_at": true,
}

// checkPeriodFields memeriksa nama dan tanggal, lalu mengembalikan
// tanggal yang sudah diparse beserta peta field yang bermasalah.
func checkPeriodFields(name, start, end string) (time.Time, time.Time, map[string]string) {
	errs := map[string]string{}

	name = strings.TrimSpace(name)
	if name == "" {
		errs["name"] = "wajib diisi"
	} else if len([]rune(name)) > maxPeriodNameLength {
		errs["name"] = "maksimal 100 karakter"
	}

	startDate, errStart := time.Parse(dateLayout, strings.TrimSpace(start))
	if errStart != nil {
		errs["start_date"] = "wajib diisi dengan format YYYY-MM-DD"
	}
	endDate, errEnd := time.Parse(dateLayout, strings.TrimSpace(end))
	if errEnd != nil {
		errs["end_date"] = "wajib diisi dengan format YYYY-MM-DD"
	}
	if errStart == nil && errEnd == nil && endDate.Before(startDate) {
		errs["end_date"] = "tidak boleh sebelum start_date"
	}

	return startDate, endDate, errs
}

// checkPeriodStatus memeriksa status baru terhadap status saat ini.
func checkPeriodStatus(current, next string) (string, bool) {
	if !validPeriodStatus[next] {
		return "harus salah satu dari draft, open, closed", false
	}
	if !CanTransitionPeriod(current, next) {
		return fmt.Sprintf("tidak bisa berpindah dari %s ke %s", current, next), false
	}
	return "", true
}

// ValidateCreatePeriod memeriksa permintaan POST dan membentuk entitas baru.
// Periode baru selalu berstatus draft. Peta error kosong berarti lolos.
func ValidateCreatePeriod(
	req model.CreatePeriodRequest,
) (model.RecruitmentPeriod, map[string]string) {
	start, end, errs := checkPeriodFields(req.Name, req.StartDate, req.EndDate)
	if len(errs) > 0 {
		return model.RecruitmentPeriod{}, errs
	}
	return model.RecruitmentPeriod{
		Name:      strings.TrimSpace(req.Name),
		StartDate: start,
		EndDate:   end,
		Status:    model.PeriodDraft,
	}, errs
}

// ValidateReplacePeriod memeriksa permintaan PUT terhadap data saat ini.
// Seluruh field wajib ada karena PUT mengganti isi secara keseluruhan.
func ValidateReplacePeriod(
	current model.RecruitmentPeriod, req model.ReplacePeriodRequest,
) (model.RecruitmentPeriod, map[string]string) {
	start, end, errs := checkPeriodFields(req.Name, req.StartDate, req.EndDate)
	if msg, ok := checkPeriodStatus(current.Status, req.Status); !ok {
		errs["status"] = msg
	}
	if len(errs) > 0 {
		return current, errs
	}

	current.Name = strings.TrimSpace(req.Name)
	current.StartDate = start
	current.EndDate = end
	current.Status = req.Status
	return current, errs
}

// ApplyPatchPeriod menyalin field yang dikirim ke data yang sudah ada.
// Field yang bernilai nil dibiarkan apa adanya.
func ApplyPatchPeriod(
	current model.RecruitmentPeriod, req model.PatchPeriodRequest,
) (model.RecruitmentPeriod, map[string]string) {
	errs := map[string]string{}
	updated := current

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		switch {
		case name == "":
			errs["name"] = "tidak boleh kosong"
		case len([]rune(name)) > maxPeriodNameLength:
			errs["name"] = "maksimal 100 karakter"
		default:
			updated.Name = name
		}
	}

	if req.StartDate != nil {
		if d, err := time.Parse(dateLayout, strings.TrimSpace(*req.StartDate)); err != nil {
			errs["start_date"] = "format harus YYYY-MM-DD"
		} else {
			updated.StartDate = d
		}
	}

	if req.EndDate != nil {
		if d, err := time.Parse(dateLayout, strings.TrimSpace(*req.EndDate)); err != nil {
			errs["end_date"] = "format harus YYYY-MM-DD"
		} else {
			updated.EndDate = d
		}
	}

	_, startBad := errs["start_date"]
	_, endBad := errs["end_date"]
	if !startBad && !endBad && updated.EndDate.Before(updated.StartDate) {
		errs["end_date"] = "tidak boleh sebelum start_date"
	}

	if req.Status != nil {
		if msg, ok := checkPeriodStatus(current.Status, *req.Status); !ok {
			errs["status"] = msg
		} else {
			updated.Status = *req.Status
		}
	}

	return updated, errs
}

// IsEmptyPatchPeriod menandai permintaan PATCH yang tidak mengubah apa pun.
func IsEmptyPatchPeriod(req model.PatchPeriodRequest) bool {
	return req.Name == nil && req.StartDate == nil && req.EndDate == nil && req.Status == nil
}
