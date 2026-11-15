// Package reportes arma la información de consulta para la gerencia y
// contabilidad. Solo lee datos de los demás paquetes; no modifica nada.
package reportes

import (
	"fmt"
	"sort"
	"strings"

	"highlights/devoluciones"
	"highlights/entregas"
	"highlights/inventario"
	"highlights/liquidacion"
	"highlights/proformas"
)

// StockBajo lista los productos cuyo stock es menor al mínimo indicado.
func StockBajo(minimo int) []string {
	var lineas []string
	for _, codigo := range inventario.ListarCodigos() {
		descripcion, _, stock, _ := inventario.ConsultarProducto(codigo)
		if stock < minimo {
			lineas = append(lineas, fmt.Sprintf("%s | %s | stock: %d", codigo, descripcion, stock))
		}
	}
	return lineas
}

// EstadoProyecto muestra, por producto, lo cotizado, entregado, devuelto y neto.
func EstadoProyecto(numeroProforma string) (string, error) {
	ruc, proyecto, estado, err := proformas.ConsultarProforma(numeroProforma)
	if err != nil {
		return "", err
	}
	if liquidacion.EstaLiquidado(numeroProforma) {
		estado = "liquidada"
	}
	items := proformas.Items(numeroProforma)
	codigos := make([]string, 0, len(items))
	for c := range items {
		codigos = append(codigos, c)
	}
	sort.Strings(codigos)

	var b strings.Builder
	fmt.Fprintf(&b, "%s - %s (cliente %s) - estado: %s\n", numeroProforma, proyecto, ruc, estado)
	fmt.Fprintf(&b, "%-10s %8s %8s %8s %8s\n", "Código", "Cotiz.", "Entreg.", "Devuel.", "Neto")
	for _, c := range codigos {
		fmt.Fprintf(&b, "%-10s %8d %8d %8d %8d\n", c, items[c],
			entregas.CantidadEntregada(numeroProforma, c),
			devoluciones.CantidadDevuelta(numeroProforma, c),
			liquidacion.CantidadNeta(numeroProforma, c))
	}
	return b.String(), nil
}

// ResumenFacturacion compara el total cotizado con el liquidado y entrega
// el valor final que contabilidad debe facturar.
func ResumenFacturacion(numeroProforma string) (string, error) {
	_, _, cotizado, err := proformas.CalcularTotal(numeroProforma, nil)
	if err != nil {
		return "", err
	}
	subtotal, iva, total, err := liquidacion.CalcularLiquidacion(numeroProforma)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Cotizado: $%.2f | Subtotal: $%.2f | IVA 15%%: $%.2f | TOTAL A FACTURAR: $%.2f | Diferencia: $%.2f",
		cotizado, subtotal, iva, total, cotizado-total), nil
}
