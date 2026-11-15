# Sistema de Gestión de Proyectos High Lights

Proyecto de la materia **Programación Orientada a Objetos** — Ingeniería en Inteligencia Artificial, UIDE.
Lenguaje: **Go**. Entrega de la **Unidad 1: funciones y paquetes** (todavía sin structs).

## ¿Qué problema resuelve?

High Lights ejecuta proyectos de iluminación: cotiza, despacha luminarias a la obra, recibe sobrantes
y al final factura. Hoy ese control se lleva en hojas sueltas, y es fácil facturar lo cotizado en lugar
de lo que realmente quedó instalado. Este sistema une en un solo flujo:

**inventario → usuarios → clientes → proformas → entregas → devoluciones → liquidación final → reportes**

y deja el valor neto del proyecto listo para facturación.

## Estructura del repositorio

```
highlights/
├── go.mod                          ← un solo módulo de Go
├── README.md
├── main.go                         ← package main: solo arranca el programa
├── inventario/inventario.go        ← package inventario
├── usuarios/usuarios.go            ← package usuarios
├── clientes/clientes.go            ← package clientes
├── proformas/proformas.go          ← package proformas
├── entregas/entregas.go            ← package entregas
├── devoluciones/devoluciones.go    ← package devoluciones
├── liquidacion/liquidacion.go      ← package liquidacion
└── reportes/reportes.go            ← package reportes
```

## Dependencias entre paquetes (sin ciclos)

```
main ──► todos los paquetes
proformas    ──► clientes, inventario
entregas     ──► proformas, inventario
devoluciones ──► entregas, inventario
liquidacion  ──► proformas, entregas, devoluciones
reportes     ──► inventario, proformas, entregas, devoluciones, liquidacion
inventario, usuarios, clientes ──► (solo biblioteca estándar)
```

Solo se usa la **biblioteca estándar** de Go (`fmt`, `errors`, `sort`, `strings`, `log`).
No hay paquetes de terceros.

## Cómo ejecutarlo

Requisito: Go 1.22 o superior.

```bash
go build ./...    # compila todos los paquetes
go run .          # ejecuta el flujo de ejemplo
```

Salida esperada (resumen):

```
PF-0001 - Oficinas Torre Norte (cliente 1790012345001) - estado: liquidada
Código       Cotiz.  Entreg.  Devuel.     Neto
DL-12W          120      120       10      110
TL-40W           30       30        0       30
Cotizado: $4002.00 | Subtotal: $3295.00 | IVA 15%: $494.25 | TOTAL A FACTURAR: $3789.25 | Diferencia: $212.75

Productos con stock bajo (< 60):
  TL-40W | Lineal LED 40W 4000K | stock: 50
```

## Alcance de esta entrega

- 8 paquetes (uno por módulo) con 29 funciones exportadas, de 3 a 6 por paquete.
- Datos en memoria (mapas y slices); se pierden al cerrar el programa.
- **No incluye** todavía: structs, base de datos, interfaz gráfica ni emisión de factura electrónica SRI.

Autora: Ana Caro
