// Package devoluciones registra la mercadería que regresa de la obra
// (sobrantes, cambios o producto defectuoso).
package devoluciones

import (
	"fmt"

	"highlights/entregas"
	"highlights/inventario"
)

// Motivos válidos de devolución.
const (
	MotivoSobrante   = "sobrante"
	MotivoCambio     = "cambio"
	MotivoDefectuoso = "defectuoso"
)

var (
	contador = 0
	devuelto = map[string]map[string]int{} // proforma, código, unidades
	notas    = map[string][]string{}       // proforma, notas de devolución
)

// RegistrarDevolucion recibe unidades desde la obra y devuelve el número de
// nota de entrega. Solo se devuelve lo entregado; lo defectuoso no vuelve al stock.
func RegistrarDevolucion(numeroProforma, codigo string, unidades int, motivo string) (string, error) {
	if !validarMotivo(motivo) {
		return "", fmt.Errorf("devoluciones: motivo %q no válido", motivo)
	}
	enObra := entregas.CantidadEntregada(numeroProforma, codigo) - CantidadDevuelta(numeroProforma, codigo)
	if unidades <= 0 || unidades > enObra {
		return "", fmt.Errorf("devoluciones: %d unidades de %s exceden lo que está en obra (%d)", unidades, codigo, enObra)
	}
	if motivo != MotivoDefectuoso {
		if err := inventario.MoverStock(codigo, unidades); err != nil {
			return "", err
		}
	}
	if devuelto[numeroProforma] == nil {
		devuelto[numeroProforma] = map[string]int{}
	}
	devuelto[numeroProforma][codigo] += unidades
	contador++
	nota := fmt.Sprintf("DV-%04d", contador)
	notas[numeroProforma] = append(notas[numeroProforma],
		fmt.Sprintf("%s | %s x%d | motivo: %s", nota, codigo, unidades, motivo))
	return nota, nil
}

// CantidadDevuelta devuelve el total de unidades devueltas de un producto.
func CantidadDevuelta(numeroProforma, codigo string) int {
	return devuelto[numeroProforma][codigo]
}

// ListarDevoluciones devuelve las notas de devolución de una proforma.
func ListarDevoluciones(numeroProforma string) []string {
	return notas[numeroProforma]
}

// validarMotivo es privada: acepta sobrante, cambio o defectuoso.
func validarMotivo(motivo string) bool {
	return motivo == MotivoSobrante || motivo == MotivoCambio || motivo == MotivoDefectuoso
}
