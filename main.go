// Programa principal del Sistema de Gestión de Proyectos High Lights.
// main solo arranca el programa y llama a los paquetes; la lógica vive en ellos.
package main

import (
	"fmt"
	"log"

	"highlights/clientes"
	"highlights/devoluciones"
	"highlights/entregas"
	"highlights/inventario"
	"highlights/liquidacion"
	"highlights/proformas"
	"highlights/reportes"
	"highlights/usuarios"
)

func main() {
	fmt.Println("=== Sistema de Gestión High Lights ===")

	revisar(usuarios.RegistrarUsuario("0102030405", "Ana Caro", usuarios.RolOperaciones))
	revisar(inventario.RegistrarProducto("DL-12W", "Downlight LED 12W 3000K", 18.50, 200))
	revisar(inventario.RegistrarProducto("TL-40W", "Lineal LED 40W 4000K", 42.00, 80))
	revisar(clientes.RegistrarCliente("1790012345001", "Constructora Andina S.A.", "facturas@andina.ec"))

	pf, err := proformas.CrearProforma("1790012345001", "Oficinas Torre Norte")
	revisar(err)
	revisar(proformas.AgregarItem(pf, "DL-12W", 120))
	revisar(proformas.AgregarItem(pf, "TL-40W", 30))
	revisar(proformas.AprobarProforma(pf))

	_, err = entregas.RegistrarEntrega(pf, "DL-12W", 120, "Residente de obra")
	revisar(err)
	_, err = entregas.RegistrarEntrega(pf, "TL-40W", 30, "Residente de obra")
	revisar(err)
	_, err = devoluciones.RegistrarDevolucion(pf, "DL-12W", 10, devoluciones.MotivoSobrante)
	revisar(err)

	if usuarios.TienePermiso("0102030405", "liquidacion") {
		revisar(liquidacion.CerrarProyecto(pf))
	}

	estado, err := reportes.EstadoProyecto(pf)
	revisar(err)
	fmt.Print(estado)
	resumen, err := reportes.ResumenFacturacion(pf)
	revisar(err)
	fmt.Println(resumen)

	fmt.Println("\nProductos con stock bajo (< 60):")
	for _, linea := range reportes.StockBajo(60) {
		fmt.Println(" ", linea)
	}
}

// revisar detiene el programa si alguna operación devolvió error.
func revisar(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
