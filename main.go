// EJERCICIO: analizador de logs.
//
// Un solo programa que atraviesa casi todo el Tour of Go. Siete etapas,
// cada una construye sobre la anterior. Hacelas en orden.
//
//   go mod init loganalyzer && go run . testdata/app.log
//
// Tiempo: 3-4 horas si venís oxidado. No es una sentada.
// Regla: terminá cada etapa y hacé que compile antes de pasar a la siguiente.
//
// ===========================================================================
// COBERTURA DEL TOUR
// ===========================================================================
//
//  Etapa 1  packages · exported names · funciones · retornos múltiples ·
//           retornos con nombre · variables · zero values · conversión de
//           tipos · inferencia · constantes · iota · Stringer
//  Etapa 2  if con statement corto · switch · switch sin condición ·
//           defer · defer apilados · errores · tipos de error
//  Etapa 3  arrays · slices · len/cap · nil slices · make · append ·
//           slices de slices · range · maps · map literals · mutar maps
//  Etapa 4  punteros · structs · punteros a structs · struct literals ·
//           valores función · closures
//  Etapa 5  métodos · receptores puntero · interfaces · implementación
//           implícita · valores de interfaz · interfaz nil · interfaz vacía ·
//           type assertions · type switches · io.Reader
//  Etapa 6  goroutines · channels · channels con buffer · range y close ·
//           select · select con default · sync.Mutex · WaitGroup
//  Etapa 7  (bonus) generics: type parameters
//
// ===========================================================================

package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ===========================================================================
// ETAPA 1 — Tipos, constantes y Stringer
// ===========================================================================
//
// Formato de una línea:
//   2026-03-14T08:23:11Z ERROR db connection refused host=pg-01 retries=3
//   └─ timestamp        └─ nivel
//                             └─ componente
//                                 └─ mensaje libre
//                                                       └─ pares clave=valor
//
// TODO 1.1: definí `type Level int` con constantes usando iota, en orden:
//           Debug, Info, Warn, Error, Fatal.
//           Que Level implemente fmt.Stringer devolviendo "DEBUG", "INFO", etc.
//           Para un nivel desconocido: "LEVEL(7)" con el número.
//
// TODO 1.2: `func ParseLevel(s string) (Level, error)` — el retorno múltiple
//           clásico. Case-insensitive.
//
// TODO 1.3: `type Entry struct` con: Timestamp time.Time, Level Level,
//           Component string, Message string, Fields map[string]string.
//           Que Entry también implemente Stringer.
//
// TODO 1.4: `func (e Entry) DurationMS() (ms int, ok bool)` — usá retornos
//           CON NOMBRE y un `return` desnudo al menos una vez, para que veas
//           cómo se lee. (Después decidí si te gusta. A mucha gente no.)
//
// Ojo con los zero values: ¿qué es un Entry vacío? ¿Y su Fields nil?
// Escribí a mano qué pasa si hacés `len(e.Fields)` sobre un map nil, y
// después comprobalo. Es distinto de escribir en él.

// ===========================================================================
// ETAPA 2 — Parseo, control de flujo y errores
// ===========================================================================
//
// TODO 2.1: `type ParseError struct { Line int; Raw string; Err error }`
//           con método Error() y Unwrap(). Errores centinela:
//           ErrMalformed, ErrUnknownLevel.
//
// TODO 2.2: `func ParseLine(n int, raw string) (Entry, error)`
//           - `if` con statement corto: `if parts := strings.Fields(raw); len(parts) < 4 {`
//           - `switch` para clasificar el nivel
//           - `switch` SIN condición (switch { case x > 10: ... }) para
//             clasificar la duración en "rápido/normal/lento". Es el reemplazo
//             idiomático de la cadena if/else-if.
//
// TODO 2.3: `func ParseFile(path string) ([]Entry, []error)` — devuelve las
//           entradas válidas Y los errores, no se detiene en el primero.
//           Usá `defer f.Close()`.
//           Meté DOS defer apilados e imprimí algo en cada uno: tenés que
//           ver con tus ojos que salen en orden inverso (LIFO).
//
// El archivo de prueba tiene dos líneas rotas a propósito. Tu parser tiene
// que sobrevivirlas y reportarlas, no morirse.

// ===========================================================================
// ETAPA 3 — Colecciones
// ===========================================================================
//
// TODO 3.1: `func CountByLevel(es []Entry) map[Level]int`
//
// TODO 3.2: `func HourlyHistogram(es []Entry) [24]int`
//           ARRAY, no slice. Fijate que se copia al pasarlo y al devolverlo,
//           a diferencia del slice. Ese contraste es el punto del TODO.
//
// TODO 3.3: `func GroupByComponent(es []Entry) map[string][]Entry`
//           Acá vive el truco favorito de Go: `m[k] = append(m[k], v)`
//           funciona aunque la clave no exista, porque append sobre un
//           slice nil hace lo correcto. Entendé POR QUÉ antes de seguir.
//
// TODO 3.4: `func Chunk(es []Entry, size int) [][]Entry`
//           Slices de slices. Cuidado: los sub-slices comparten el array
//           subyacente con el original. Escribí un comentario explicando
//           cuándo eso te puede morder.
//
// TODO 3.5: imprimí len y cap de un slice mientras le hacés append en un
//           loop de 20 iteraciones. Mirá cómo crece cap. No es lineal.

// ===========================================================================
// ETAPA 4 — Punteros, structs y closures
// ===========================================================================
//
// TODO 4.1: `type Stats struct` con Total, ByLevel map[Level]int,
//           SlowestMS int, SlowestEntry *Entry.
//           El puntero es deliberado: ¿por qué *Entry y no Entry?
//
// TODO 4.2: `func (s *Stats) Add(e Entry)` — receptor PUNTERO porque muta.
//           Escribí también `func (s Stats) Summary() string` con receptor
//           VALOR porque solo lee. Que la diferencia sea una decisión tuya
//           y no una casualidad.
//
// TODO 4.3: `func Counter() func(Level) int` — closure que recuerda cuántas
//           veces vio cada nivel. Mismo espíritu que el Fibonacci del Tour.
//
// TODO 4.4: `func FilterFunc(pred func(Entry) bool) func([]Entry) []Entry`
//           Una función que devuelve una función. Usala para componer:
//           soloErrores := FilterFunc(func(e Entry) bool { return e.Level >= Error })
//
// TODO 4.5: ordená las entradas por duración usando sort.Slice, que recibe
//           un valor función. Después ordenalas por componente y desempatá
//           por timestamp.

// ===========================================================================
// ETAPA 5 — Interfaces
// ===========================================================================
//
// TODO 5.1: `type Filter interface { Match(Entry) bool }`
//           Tres implementaciones que NO declaran que implementan nada
//           (implementación implícita — es el corazón del sistema de tipos):
//             - LevelFilter (un Level mínimo)
//             - ComponentFilter (un string)
//             - AndFilter ([]Filter, todos tienen que matchear)
//
// TODO 5.2: `func Apply(es []Entry, f Filter) []Entry`
//           Probalo pasándole un Filter nil. Va a explotar. Ahora hacé que
//           un nil signifique "pasa todo" y fijate la diferencia entre una
//           interfaz nil y una interfaz con valor nil adentro. Este es el
//           tema que más confunde de Go y vale que lo toques ahora.
//
// TODO 5.3: `func Describe(v any)` con un TYPE SWITCH que maneje
//           Entry, Level, []Entry, string, int, error y default.
//           (`any` es alias de `interface{}` — la interfaz vacía.)
//
// TODO 5.4: cambiá ParseFile para que reciba un `io.Reader` en vez de una
//           ruta. Ahora podés testearlo con strings.NewReader("...") sin
//           tocar el disco. Esta única línea es la razón por la que io.Reader
//           está en el Tour, y es el patrón que más vas a usar en tu vida.

// ===========================================================================
// ETAPA 6 — Concurrencia
// ===========================================================================
//
// TODO 6.1: `func ParseConcurrent(paths []string) ([]Entry, []error)`
//           Una goroutine por archivo, resultados por channel, sync.WaitGroup
//           para esperar a todas.
//           Duplicá testdata/app.log en 4 archivos para tener con qué probar.
//
// TODO 6.2: channel CON BUFFER de tamaño len(paths). Corré también con
//           buffer 0 y explicá en un comentario qué cambia.
//
// TODO 6.3: cerrá el channel y consumilo con `for e := range ch`.
//           Preguntas que tenés que poder responder: ¿quién cierra? ¿qué
//           pasa si cerrás dos veces? ¿si escribís en uno cerrado?
//           Probá las dos y mirá el panic.
//
// TODO 6.4: un agregador con `select` sobre tres canales: entradas, errores
//           y un `done chan struct{}` para cortar.
//           Agregá un `case <-time.After(2*time.Second)` como timeout.
//
// TODO 6.5: un `select` con `default` para un "drenado no bloqueante":
//           vaciar lo que haya en el channel sin quedarse esperando.
//
// TODO 6.6: un `*Stats` compartido protegido con sync.Mutex, actualizado
//           desde todas las goroutines.
//           Después corré `go run -race .` — si no marca nada, quitá el
//           mutex a propósito y volvé a correrlo. Ver el detector de races
//           encontrarlo vale más que cualquier explicación.

// ===========================================================================
// ETAPA 7 — Bonus: generics
// ===========================================================================
//
// No lo vas a necesitar en las próximas semanas. Es para cerrar el Tour.
//
// TODO 7.1: `func Map[T, U any](in []T, f func(T) U) []U`
// TODO 7.2: `func GroupBy[T any, K comparable](in []T, key func(T) K) map[K][]T`
//           Reescribí GroupByComponent usando esta. Fijate por qué K tiene
//           que ser `comparable` y no `any`.

// ===========================================================================
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

const UTC_TIMEZONE = -5

var ErrUnknownLevel = errors.New("Nivel de log desconocido")

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return fmt.Sprintf("LEVEL(%d)", l)
	}
}

func ParseLevel(l string) (Level, error) {
	switch strings.ToUpper(l) {
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "WARN":
		return WARN, nil
	case "ERROR":
		return ERROR, nil
	case "FATAL":
		return FATAL, nil
	default:
		return -1, ErrUnknownLevel
	}
}

type Entry struct {
	Timestamp time.Time
	Level     Level
	Component string
	Message   string
	Fields    map[string]string
}

func (entry Entry) String() string {

	if !entry.Timestamp.IsZero() {
		entry.Timestamp = entry.Timestamp.Local()
	}

	return fmt.Sprintf("%s | %s | %s | %s | %s",
		entry.Timestamp.Format("15:04:05"),
		entry.Level.String(),
		entry.Component,
		entry.Message,
		entry.Fields["dur_ms"])
}

// TODO 1.3: `type Entry struct` con: Timestamp time.Time, Level Level,
//           Component string, Message string, Fields map[string]string.
//           Que Entry también implemente Stringer.

func main() {
	fmt.Println("empezá por la etapa 1")
}

// ===========================================================================
// CHECKLIST FINAL
// ===========================================================================
//
//   [ ] Compila sin warnings: go vet ./...
//   [ ] staticcheck ./... sale limpio
//   [ ] go run -race . no reporta nada
//   [ ] Ningún error ignorado con _
//   [ ] Cada receptor puntero vs valor es una decisión, no un default
//   [ ] Podés explicar en voz alta: interfaz nil vs interfaz con valor nil
//   [ ] Podés explicar por qué append sobre un slice nil funciona
//   [ ] Sobrevive las dos líneas rotas de testdata/app.log
//
// Si algo de la lista no te sale, ese es exactamente el tema para repasar.
