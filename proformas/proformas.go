// Package proformas gestiona las cotizaciones de cada proyecto:
// ítems, precios, IVA y aprobación por parte del cliente.
package proformas

import (
	"errors"
	"fmt"

	"highlights/clientes"
	"highlights/inventario"
)

// IVA vigente en Ecuador.
const IVA = 0.15

// Estados posibles de una proforma.
const (
	EstadoBorrador = "borrador"
	EstadoAprobada = "aprobada"
)

var (
	contador = 0
	cliente  = map[string]string{}             // número RUC
	proyecto = map[string]string{}             // número nombre del proyecto
	estado   = map[string]string{}             // número estado
	cantidad = map[string]map[string]int{}     // número código unidades
	precio   = map[string]map[string]float64{} // número código precio congelado
)

// CrearProforma abre una proforma en borrador y devuelve su número.
func CrearProforma(rucCliente, nombreProyecto string) (string, error) {
	if !clientes.Existe(rucCliente) {
		return "", fmt.Errorf("proformas: cliente %s no registrado", rucCliente)
	}
	if nombreProyecto == "" {
		return "", errors.New("proformas: el nombre del proyecto es obligatorio")
	}
	contador++
	numero := fmt.Sprintf("PF-%04d", contador)
	cliente[numero] = rucCliente
	proyecto[numero] = nombreProyecto
	estado[numero] = EstadoBorrador
	cantidad[numero] = map[string]int{}
	precio[numero] = map[string]float64{}
	return numero, nil
}

// AgregarItem añade un producto a una proforma en borrador y congela su precio.
func AgregarItem(numero, codigo string, unidades int) error {
	if estado[numero] != EstadoBorrador {
		return fmt.Errorf("proformas: %s no existe o ya fue aprobada", numero)
	}
	if unidades <= 0 {
		return errors.New("proformas: la cantidad debe ser > 0")
	}
	_, p, _, err := inventario.ConsultarProducto(codigo)
	if err != nil {
		return err
	}
	cantidad[numero][codigo] += unidades
	precio[numero][codigo] = p
	return nil
}

// AprobarProforma marca la proforma como aprobada. Solo así puede despacharse.
func AprobarProforma(numero string) error {
	if estado[numero] != EstadoBorrador {
		return fmt.Errorf("proformas: %s no existe o ya fue aprobada", numero)
	}
	if len(cantidad[numero]) == 0 {
		return errors.New("proformas: no se puede aprobar una proforma sin ítems")
	}
	estado[numero] = EstadoAprobada
	return nil
}

// ConsultarProforma devuelve el RUC del cliente, el proyecto y el estado.
func ConsultarProforma(numero string) (string, string, string, error) {
	if _, ok := estado[numero]; !ok {
		return "", "", "", fmt.Errorf("proformas: %s no existe", numero)
	}
	return cliente[numero], proyecto[numero], estado[numero], nil
}

// Items devuelve una copia de los ítems cotizados: código unidades.
func Items(numero string) map[string]int {
	copia := map[string]int{}
	for c, u := range cantidad[numero] {
		copia[c] = u
	}
	return copia
}

// CalcularTotal valora la proforma y devuelve subtotal, IVA y total.
// Si unidades es nil usa las cantidades cotizadas. Si no, usa las que reciba
// así liquidacion puede valorar las cantidades realmente instaladas.
func CalcularTotal(numero string, unidades map[string]int) (float64, float64, float64, error) {
	if _, ok := estado[numero]; !ok {
		return 0, 0, 0, fmt.Errorf("proformas: %s no existe", numero)
	}
	if unidades == nil {
		unidades = cantidad[numero]
	}
	subtotal := 0.0
	for codigo, u := range unidades {
		subtotal += float64(u) * precio[numero][codigo]
	}
	iva := subtotal * IVA
	return subtotal, iva, subtotal + iva, nil
}
