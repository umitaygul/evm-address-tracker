package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/umitaygul/evm-address-tracker/internal/repository"
)

type AddressHandler struct {
	addresses *repository.AddressRepository
}

func NewAddressHandler(addresses *repository.AddressRepository) *AddressHandler {
	return &AddressHandler{addresses: addresses}
}

type addAddressRequest struct {
	ChainID int64  `json:"chain_id" binding:"required"`
	Address string `json:"address" binding:"required"`
}

func (h *AddressHandler) Add(c *gin.Context) {
	var req addAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isValidEVMAddress(req.Address) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid EVM address"})
		return
	}

	userID := c.GetString("user_id")

	addr, err := h.addresses.Create(c.Request.Context(), userID, req.ChainID, strings.ToLower(req.Address))
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			c.JSON(http.StatusConflict, gin.H{"error": "address already watched on this chain"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, addr)
}

func (h *AddressHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")

	addresses, err := h.addresses.ListByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"addresses": addresses})
}

func (h *AddressHandler) Get(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	addr, err := h.addresses.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "address not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, addr)
}

func (h *AddressHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	err := h.addresses.Delete(c.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "address not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// isValidEVMAddress checks for a basic 0x-prefixed 20-byte hex address.
func isValidEVMAddress(addr string) bool {
	if len(addr) != 42 {
		return false
	}
	if addr[:2] != "0x" && addr[:2] != "0X" {
		return false
	}
	for _, c := range addr[2:] {
		if !isHexChar(c) {
			return false
		}
	}
	return true
}

func isHexChar(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}
