// Package usuarios gestiona al personal interno que usa el sistema
// y los permisos que tiene cada rol.
package usuarios

import (
	"errors"
	"fmt"
	"sort"
)

// Roles válidos del sistema.
const (
	RolAdministrador = "administrador"
	RolVentas        = "ventas"
	RolBodega        = "bodega"
	RolOperaciones   = "operaciones"
)

var (
	nombres = map[string]string{} // cédula, nombre
	roles   = map[string]string{} // cédula, rol

	// permisos indica qué módulos puede usar cada rol.
	permisos = map[string][]string{
		RolAdministrador: {"inventario", "usuarios", "clientes", "proformas", "entregas", "devoluciones", "liquidacion", "reportes"},
		RolVentas:        {"clientes", "proformas", "reportes"},
		RolBodega:        {"inventario", "entregas", "devoluciones"},
		RolOperaciones:   {"entregas", "devoluciones", "liquidacion", "reportes"},
	}
)

// RegistrarUsuario crea un usuario con un rol válido.
func RegistrarUsuario(cedula, nombre, rol string) error {
	if len(cedula) != 10 || nombre == "" {
		return errors.New("usuarios: cédula de 10 dígitos y nombre son obligatorios")
	}
	if !validarRol(rol) {
		return fmt.Errorf("usuarios: rol %q no válido", rol)
	}
	if _, ok := nombres[cedula]; ok {
		return fmt.Errorf("usuarios: la cédula %s ya está registrada", cedula)
	}
	nombres[cedula] = nombre
	roles[cedula] = rol
	return nil
}

// TienePermiso indica si un usuario puede usar un módulo según su rol.
func TienePermiso(cedula, modulo string) bool {
	for _, m := range permisos[roles[cedula]] {
		if m == modulo {
			return true
		}
	}
	return false
}

// ListarUsuarios devuelve una línea por usuario, ordenada por cédula.
func ListarUsuarios() []string {
	cedulas := make([]string, 0, len(nombres))
	for c := range nombres {
		cedulas = append(cedulas, c)
	}
	sort.Strings(cedulas)
	lineas := make([]string, 0, len(cedulas))
	for _, c := range cedulas {
		lineas = append(lineas, fmt.Sprintf("%s | %s | %s", c, nombres[c], roles[c]))
	}
	return lineas
}

// validarRol es privada: solo la usa este paquete.
func validarRol(rol string) bool {
	_, ok := permisos[rol]
	return ok
}
