package parseFormData

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func ParseFormData(r *http.Request) (*Event, error) {
	const op = "handlers.event.parseFormData"
	err := r.ParseMultipartForm(MaxSizeForm)
	mForm := r.MultipartForm
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var event Event

	if mForm.Value["price"] != nil {
		if event.Price, err = strconv.ParseUint(mForm.Value["price"][0], 10, 64); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}
	if mForm.Value["restrictions"] != nil {
		if event.Restrictions, err = strconv.ParseUint(mForm.Value["restrictions"][0], 10, 64); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	if mForm.Value["date"] != nil {
		if event.Date, err = time.Parse(time.RFC3339, mForm.Value["date"][0]); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	if mForm.Value["feature"] != nil {
		event.Feature = mForm.Value["feature"]
	}
	if mForm.Value["city"] != nil {
		event.City = mForm.Value["city"][0]
	}
	if mForm.Value["address"] != nil {
		event.Address = mForm.Value["address"][0]
	}
	if mForm.Value["name"] != nil {
		event.Name = mForm.Value["name"][0]
	}
	if mForm.Value["description"] != nil {
		event.Description = mForm.Value["description"][0]
	}

	files, ok := mForm.File["img_path"]
	ext := strings.Split(files[0].Filename, ".")
	path := fmt.Sprintf("%s/%s.%s", ImageFolder, uuid.NewString(), ext[len(ext)-1])
	if ok {
		img, err := files[0].Open()
		if err != nil {
			slog.Error("couldn't open sent image", slogResponse.SlogOp(op), slogResponse.SlogErr(err))
			return &event, fmt.Errorf("%s: %w", op, err)
		}
		localFile, err := os.Create(path)
		if err != nil {
			slog.Error("couldn't create image copy", slogResponse.SlogOp(op), slogResponse.SlogErr(err))
			return &event, fmt.Errorf("%s: %w", op, err)
		}
		if _, err = io.Copy(localFile, img); err != nil {
			slog.Error("couldn't save image", slogResponse.SlogOp(op), slogResponse.SlogErr(err))
			return &event, fmt.Errorf("%s: %w", op, err)
		}
	}

	event.ImgPath = path

	return &event, nil
}
