// Package clientes gestiona las empresas a quienes High Lights
// cotiza, entrega y factura proyectos de iluminación.
package clientes

import (
	"errors"
	"fmt"
	"strings"
)

var (
	razonesSociales = map[string]string{} // RUC -> razón social
	correos         = map[string]string{} // RUC -> correo de facturación
)

// RegistrarCliente agrega un cliente nuevo identificado por su RUC.
func RegistrarCliente(ruc, razonSocial, correo string) error {
	if !validarRUC(ruc) {
		return fmt.Errorf("clientes: RUC %q no válido (13 dígitos terminados en 001)", ruc)
	}
	if razonSocial == "" || !strings.Contains(correo, "@") {
		return errors.New("clientes: razón social y correo válido son obligatorios")
	}
	if Existe(ruc) {
		return fmt.Errorf("clientes: el RUC %s ya está registrado", ruc)
	}
	razonesSociales[ruc] = razonSocial
	correos[ruc] = correo
	return nil
}

// Existe indica si el RUC ya está registrado.
func Existe(ruc string) bool {
	_, ok := razonesSociales[ruc]
	return ok
}

// BuscarCliente devuelve la razón social y el correo de facturación.
func BuscarCliente(ruc string) (string, string, error) {
	if !Existe(ruc) {
		return "", "", fmt.Errorf("clientes: RUC %s no encontrado", ruc)
	}
	return razonesSociales[ruc], correos[ruc], nil
}

// validarRUC hace una validación básica del RUC ecuatoriano:
// 13 dígitos numéricos que terminan en "001".
func validarRUC(ruc string) bool {
	if len(ruc) != 13 || !strings.HasSuffix(ruc, "001") {
		return false
	}
	for _, r := range ruc {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
