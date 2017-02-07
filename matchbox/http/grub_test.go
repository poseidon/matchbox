package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"context"
	logtest "github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	fake "github.com/coreos/matchbox/matchbox/storage/testfakes"
)

func TestGrubHandler(t *testing.T) {
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	h := srv.grubHandler()
	ctx := withProfile(context.Background(), fake.Profile)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(ctx, w, req)
	// assert that:
	// - the Profile's NetBoot config is rendered as a GRUB2 config
	expectedScript := `default=0
fallback=1
timeout=1
menuentry "CoreOS (EFI)" {
echo "Loading kernel"
linuxefi "/image/kernel" a=b c
echo "Loading initrd"
initrdefi  "/image/initrd_a" "/image/initrd_b"
}
menuentry "CoreOS (BIOS)" {
echo "Loading kernel"
linux "/image/kernel" a=b c
echo "Loading initrd"
initrd  "/image/initrd_a" "/image/initrd_b"
}
`
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedScript, w.Body.String())
}
