// Package inventario gestiona los productos de bodega de High Lights:
// luminarias y accesorios, con su precio y existencias (stock).
package inventario

import (
	"errors"
	"fmt"
	"sort"
)

// Datos internos del paquete (minúscula = privados).

var (
	descripciones = map[string]string{}
	precios       = map[string]float64{}
	stock         = map[string]int{}
)

// RegistrarProducto agrega un producto nuevo al inventario.
func RegistrarProducto(codigo, descripcion string, precio float64, cantidad int) error {
	if codigo == "" || descripcion == "" {
		return errors.New("inventario: código y descripción son obligatorios")
	}
	if precio <= 0 || cantidad < 0 {
		return errors.New("inventario: precio debe ser > 0 y cantidad >= 0")
	}
	if _, ok := descripciones[codigo]; ok {
		return fmt.Errorf("inventario: el producto %s ya existe", codigo)
	}
	descripciones[codigo] = descripcion
	precios[codigo] = precio
	stock[codigo] = cantidad
	return nil
}

// ConsultarProducto devuelve la descripción, el precio sin IVA y el stock de un producto.
func ConsultarProducto(codigo string) (string, float64, int, error) {
	if _, ok := descripciones[codigo]; !ok {
		return "", 0, 0, fmt.Errorf("inventario: producto %s no encontrado", codigo)
	}
	return descripciones[codigo], precios[codigo], stock[codigo], nil
}

// MoverStock suma o resta unidades.
// Se usa cuando la mercadería sale a obra o regresa de ella.
func MoverStock(codigo string, cantidad int) error {
	if _, ok := descripciones[codigo]; !ok {
		return fmt.Errorf("inventario: producto %s no encontrado", codigo)
	}
	if stock[codigo]+cantidad < 0 {
		return fmt.Errorf("inventario: stock insuficiente de %s (disponible %d)", codigo, stock[codigo])
	}
	stock[codigo] += cantidad
	return nil
}

// ListarCodigos devuelve los códigos de todos los productos ordenados.
func ListarCodigos() []string {
	codigos := make([]string, 0, len(descripciones))
	for c := range descripciones {
		codigos = append(codigos, c)
	}
	sort.Strings(codigos)
	return codigos
}
