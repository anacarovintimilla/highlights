// Package liquidacion cruza lo cotizado, lo entregado y lo devuelto para
// obtener el valor real del proyecto y dejarlo listo para facturación.
package liquidacion

import (
	"fmt"

	"highlights/devoluciones"
	"highlights/entregas"
	"highlights/proformas"
)

var cerrados = map[string]bool{} // proforma liquidada

// CantidadNeta devuelve las unidades realmente instaladas: entregadas menos devueltas.
func CantidadNeta(numeroProforma, codigo string) int {
	return entregas.CantidadEntregada(numeroProforma, codigo) - devoluciones.CantidadDevuelta(numeroProforma, codigo)
}

// CalcularLiquidacion devuelve subtotal, IVA y total según las cantidades netas.
func CalcularLiquidacion(numeroProforma string) (float64, float64, float64, error) {
	netas := map[string]int{}
	for codigo := range proformas.Items(numeroProforma) {
		netas[codigo] = CantidadNeta(numeroProforma, codigo)
	}
	return proformas.CalcularTotal(numeroProforma, netas)
}

// CerrarProyecto marca la proforma como liquidada, luego pasa a facturación.
func CerrarProyecto(numeroProforma string) error {
	_, _, estado, err := proformas.ConsultarProforma(numeroProforma)
	if err != nil {
		return err
	}
	if estado != proformas.EstadoAprobada {
		return fmt.Errorf("liquidacion: la proforma %s no está aprobada", numeroProforma)
	}
	if cerrados[numeroProforma] {
		return fmt.Errorf("liquidacion: la proforma %s ya fue liquidada", numeroProforma)
	}
	cerrados[numeroProforma] = true
	return nil
}

// EstaLiquidado indica si el proyecto ya fue cerrado.
func EstaLiquidado(numeroProforma string) bool {
	return cerrados[numeroProforma]
}
