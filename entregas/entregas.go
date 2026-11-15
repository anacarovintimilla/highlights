// Package entregas registra la salida de mercadería desde bodega hacia la obra,
// mediante guías de remisión ligadas a una proforma aprobada.
package entregas

import (
	"errors"
	"fmt"

	"highlights/inventario"
	"highlights/proformas"
)

var (
	contador  = 0
	entregado = map[string]map[string]int{} // proforma -> código -> unidades
	guias     = map[string][]string{}       // proforma -> guías emitidas
)

// RegistrarEntrega despacha unidades a la obra, descuenta el stock y devuelve
// el número de guía. No permite entregar más de lo cotizado.
func RegistrarEntrega(numeroProforma, codigo string, unidades int, responsable string) (string, error) {
	_, _, estado, err := proformas.ConsultarProforma(numeroProforma)
	if err != nil {
		return "", err
	}
	if estado != proformas.EstadoAprobada {
		return "", fmt.Errorf("entregas: la proforma %s no está aprobada", numeroProforma)
	}
	if responsable == "" {
		return "", errors.New("entregas: indique quién recibe en obra")
	}
	if unidades <= 0 || unidades > pendiente(numeroProforma, codigo) {
		return "", fmt.Errorf("entregas: %d unidades de %s exceden lo pendiente (%d)",
			unidades, codigo, pendiente(numeroProforma, codigo))
	}
	if err := inventario.MoverStock(codigo, -unidades); err != nil {
		return "", err
	}
	if entregado[numeroProforma] == nil {
		entregado[numeroProforma] = map[string]int{}
	}
	entregado[numeroProforma][codigo] += unidades
	contador++
	guia := fmt.Sprintf("GR-%04d", contador)
	guias[numeroProforma] = append(guias[numeroProforma],
		fmt.Sprintf("%s | %s x%d | recibe: %s", guia, codigo, unidades, responsable))
	return guia, nil
}

// CantidadEntregada devuelve el total de unidades entregadas de un producto.
func CantidadEntregada(numeroProforma, codigo string) int {
	return entregado[numeroProforma][codigo]
}

// ListarGuias devuelve las guías de remisión emitidas para una proforma.
func ListarGuias(numeroProforma string) []string {
	return guias[numeroProforma]
}

// pendiente es privada: unidades cotizadas que aún no se despachan.
func pendiente(numeroProforma, codigo string) int {
	return proformas.Items(numeroProforma)[codigo] - CantidadEntregada(numeroProforma, codigo)
}
