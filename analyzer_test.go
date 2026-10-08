// ARNÉS DE PRUEBAS — tu termómetro para cada etapa.
//
// CÓMO SE USA
// -----------
// Todo está comentado salvo la etapa 1. A medida que avanzás, descomentás
// el bloque de la etapa que estás haciendo.
//
//   go test -run TestEtapa1 -v ./...
//
// El ciclo es siempre el mismo, y es el ciclo entero de Learn Go with Tests:
//
//   1. Descomentás el bloque de la etapa      → no compila
//   2. Definís los tipos/firmas vacíos        → compila, test ROJO
//   3. Implementás                            → test VERDE
//   4. Limpiás sabiendo que si rompés algo, el test te avisa
//
// Que no compile en el paso 1 NO es un error tuyo: en Go el compilador es
// la primera línea de tests. "No compila" es información, no fracaso.

package main

import (
	"testing"
)

// ===========================================================================
// ETAPA 1 — Tipos, constantes, Stringer          (no usa el log)
// ===========================================================================

func TestEtapa1_LevelStringer(t *testing.T) {
	casos := []struct {
		nombre string
		nivel  Level
		quiero string
	}{
		{"debug", Debug, "DEBUG"},
		{"info", Info, "INFO"},
		{"warn", Warn, "WARN"},
		{"error", Error, "ERROR"},
		{"fatal", Fatal, "FATAL"},
		{"desconocido", Level(7), "LEVEL(7)"},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := c.nivel.String(); got != c.quiero {
				t.Errorf("Level(%d).String() = %q, quiero %q", c.nivel, got, c.quiero)
			}
		})
	}
}

func TestEtapa1_LevelOrdenados(t *testing.T) {
	// El orden importa: Error >= Warn es la base de LevelFilter en la etapa 5.
	if !(Debug < Info && Info < Warn && Warn < Error && Error < Fatal) {
		t.Fatal("las constantes iota no están en orden ascendente")
	}
}

func TestEtapa1_ParseLevel(t *testing.T) {
	if got, err := ParseLevel("error"); err != nil || got != Error {
		t.Errorf(`ParseLevel("error") = %v, %v; quiero Error, nil`, got, err)
	}
	if got, err := ParseLevel("ERROR"); err != nil || got != Error {
		t.Errorf("ParseLevel debe ser case-insensitive, devolvió %v, %v", got, err)
	}
	if _, err := ParseLevel("TRENCH"); err == nil {
		t.Error("ParseLevel(\"TRENCH\") debería fallar: ese nivel no existe")
	}
}

func TestEtapa1_ZeroValues(t *testing.T) {
	var e Entry // zero value

	// len() sobre un map nil funciona. Escribirle NO.
	if len(e.Fields) != 0 {
		t.Error("un Entry vacío debería tener Fields de largo 0")
	}
	if e.Fields != nil {
		t.Error("Fields de un Entry vacío debería ser nil, no un map vacío")
	}
	// Probá a mano, fuera del test, qué pasa con: e.Fields["x"] = "y"
}

// ===========================================================================
// ETAPA 2 — Parseo y errores                     (usa líneas del log)
// ===========================================================================
/*
func TestEtapa2_ParseLineValida(t *testing.T) {
	raw := "2026-03-14T08:23:12Z ERROR db connection refused host=pg-01 retries=3"

	e, err := ParseLine(1, raw)
	if err != nil {
		t.Fatalf("no debería fallar: %v", err)
	}
	if e.Level != Error {
		t.Errorf("Level = %v, quiero ERROR", e.Level)
	}
	if e.Component != "db" {
		t.Errorf("Component = %q, quiero \"db\"", e.Component)
	}
	if e.Fields["host"] != "pg-01" {
		t.Errorf("Fields[host] = %q, quiero \"pg-01\"", e.Fields["host"])
	}
	if e.Message != "connection refused" {
		t.Errorf("Message = %q — los pares clave=valor NO van en el mensaje", e.Message)
	}
}

func TestEtapa2_ParseLineRota(t *testing.T) {
	// Estas dos líneas están en testdata/app.log a propósito.
	casos := []struct {
		nombre string
		raw    string
		quiero error
	}{
		{"malformada", "esta linea esta rota a proposito", ErrMalformed},
		{"nivel inexistente", "2026-03-14T10:00:02Z TRENCH nivel inexistente a proposito", ErrUnknownLevel},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			_, err := ParseLine(6, c.raw)
			if err == nil {
				t.Fatal("debería fallar")
			}
			// Si esto no pasa, te falta Unwrap() en ParseError.
			if !errors.Is(err, c.quiero) {
				t.Errorf("errors.Is no encuentra %v en %v — ¿implementaste Unwrap?", c.quiero, err)
			}
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatal("errors.As no extrae *ParseError")
			}
			if pe.Line != 6 {
				t.Errorf("ParseError.Line = %d, quiero 6", pe.Line)
			}
		})
	}
}

func TestEtapa2_DuracionOpcional(t *testing.T) {
	conDur, _ := ParseLine(1, "2026-03-14T08:23:11Z INFO auth login ok dur_ms=42")
	if ms, ok := conDur.DurationMS(); !ok || ms != 42 {
		t.Errorf("DurationMS() = %d, %v; quiero 42, true", ms, ok)
	}

	sinDur, _ := ParseLine(2, "2026-03-14T10:00:01Z INFO worker job started job_id=772")
	if _, ok := sinDur.DurationMS(); ok {
		t.Error("sin dur_ms, ok debería ser false")
	}
}
*/

// ===========================================================================
// ETAPA 3 — Colecciones                          (sobre el log entero)
// ===========================================================================
/*
// entradasDePrueba carga el log real. Si tu parser está bien, saltea las
// 2 líneas rotas y devuelve 13 entradas.
func entradasDePrueba(t *testing.T) []Entry {
	t.Helper()
	f, err := os.Open("testdata/app.log")
	if err != nil {
		t.Fatalf("no pude abrir el log: %v", err)
	}
	defer f.Close()

	es, errs := ParseFile(f) // en la etapa 2 puede recibir un path; ajustá
	if len(errs) != 2 {
		t.Fatalf("esperaba exactamente 2 líneas rotas, hubo %d", len(errs))
	}
	return es
}

func TestEtapa3_CountByLevel(t *testing.T) {
	es := entradasDePrueba(t)
	c := CountByLevel(es)

	// Contá a mano en testdata/app.log antes de correr esto.
	quiero := map[Level]int{Debug: 2, Info: 6, Warn: 2, Error: 3, Fatal: 1}
	for nivel, n := range quiero {
		if c[nivel] != n {
			t.Errorf("%v: conté %d, quiero %d", nivel, c[nivel], n)
		}
	}
}

func TestEtapa3_Histograma(t *testing.T) {
	es := entradasDePrueba(t)
	h := HourlyHistogram(es)

	if h[8] != 5 {
		t.Errorf("hora 08: %d entradas, quiero 5", h[8])
	}
	if h[23] != 0 {
		t.Errorf("hora 23 debería estar vacía, hay %d", h[23])
	}

	// La prueba de que es un ARRAY y no un slice: modificar la copia
	// no toca el original.
	copia := h
	copia[0] = 999
	if h[0] == 999 {
		t.Error("si esto pasa, devolviste un slice, no un array [24]int")
	}
}

func TestEtapa3_GroupByComponent(t *testing.T) {
	es := entradasDePrueba(t)
	g := GroupByComponent(es)

	if len(g["db"]) != 3 {
		t.Errorf("db: %d entradas, quiero 3", len(g["db"]))
	}
	if _, existe := g["inexistente"]; existe {
		t.Error("no deberías crear claves vacías")
	}
}

func TestEtapa3_Chunk(t *testing.T) {
	es := entradasDePrueba(t) // 13 entradas
	cs := Chunk(es, 5)

	if len(cs) != 3 {
		t.Fatalf("13 entradas en chunks de 5 = 3 grupos, devolviste %d", len(cs))
	}
	if len(cs[2]) != 3 {
		t.Errorf("el último chunk debería tener 3, tiene %d", len(cs[2]))
	}
}
*/

// ===========================================================================
// ETAPA 4 — Punteros, structs, closures
// ===========================================================================
/*
func TestEtapa4_StatsMutaPorPuntero(t *testing.T) {
	es := entradasDePrueba(t)
	s := &Stats{}
	for _, e := range es {
		s.Add(e)
	}
	if s.Total != len(es) {
		t.Errorf("Total = %d, quiero %d — ¿Add tiene receptor puntero?", s.Total, len(es))
	}
	// La entrada más lenta del log es el job 771, con 15204 ms.
	if s.SlowestMS != 15204 {
		t.Errorf("SlowestMS = %d, quiero 15204", s.SlowestMS)
	}
	if s.SlowestEntry == nil {
		t.Fatal("SlowestEntry quedó nil")
	}
}

func TestEtapa4_Closure(t *testing.T) {
	c := Counter()
	c(Error)
	c(Error)
	c(Info)
	if n := c(Error); n != 3 {
		t.Errorf("el closure contó %d errores, quiero 3 — ¿capturó el map por referencia?", n)
	}
	// Dos counters distintos NO comparten estado:
	otro := Counter()
	if n := otro(Error); n != 1 {
		t.Errorf("counters independientes: %d, quiero 1", n)
	}
}

func TestEtapa4_FilterFunc(t *testing.T) {
	es := entradasDePrueba(t)
	soloGraves := FilterFunc(func(e Entry) bool { return e.Level >= Error })

	graves := soloGraves(es)
	if len(graves) != 4 { // 3 ERROR + 1 FATAL
		t.Errorf("%d entradas graves, quiero 4", len(graves))
	}
	// Y el original no se tocó:
	if len(es) != 13 {
		t.Error("el filtro mutó el slice de entrada")
	}
}
*/

// ===========================================================================
// ETAPA 5 — Interfaces
// ===========================================================================
/*
func TestEtapa5_Filtros(t *testing.T) {
	es := entradasDePrueba(t)

	casos := []struct {
		nombre string
		filtro Filter
		quiero int
	}{
		{"nivel", LevelFilter{Min: Error}, 4},
		{"componente", ComponentFilter{Name: "db"}, 3},
		{"and", AndFilter{LevelFilter{Min: Error}, ComponentFilter{Name: "db"}}, 2},
		{"nil pasa todo", nil, 13},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := len(Apply(es, c.filtro)); got != c.quiero {
				t.Errorf("%d entradas, quiero %d", got, c.quiero)
			}
		})
	}
}

func TestEtapa5_InterfazNilVsValorNil(t *testing.T) {
	// EL test que más enseña de todo el ejercicio.
	var f Filter // interfaz nil
	if f != nil {
		t.Error("una interfaz sin asignar es nil")
	}

	var lf *LevelFilter // puntero nil
	f = lf              // interfaz CON un valor nil adentro
	if f == nil {
		t.Error("¡ojo! una interfaz con un puntero nil adentro NO es nil")
	}
	// Si este test te sorprende, releelo hasta que no te sorprenda.
	// Es la causa del bug de nil más común de Go.
}

func TestEtapa5_LeeDeCualquierReader(t *testing.T) {
	// Esta es la razón de ser de io.Reader: testear sin tocar el disco.
	r := strings.NewReader(
		"2026-03-14T08:00:00Z INFO test todo bien\n" +
			"linea rota\n")

	es, errs := ParseFile(r)
	if len(es) != 1 || len(errs) != 1 {
		t.Errorf("quiero 1 entrada y 1 error, tengo %d y %d", len(es), len(errs))
	}
}
*/

// ===========================================================================
// ETAPA 6 — Concurrencia
// ===========================================================================
//
// Acá el test cambia de naturaleza. No alcanza con verificar el resultado:
// el código concurrente incorrecto pasa los tests el 95% de las veces.
//
// Lo que de verdad prueba esta etapa es:
//
//     go test -race -count=20 -run TestEtapa6 ./...
//
//   -race   detector de condiciones de carrera
//   -count  corre el test 20 veces; las races no son deterministas
//
// Experimento obligatorio: hacé que pase en verde, después QUITÁ el mutex
// y volvé a correr con -race. Verlo fallar enseña más que cualquier lectura.
/*
func TestEtapa6_ParseConcurrente(t *testing.T) {
	paths := []string{
		"testdata/app.log", "testdata/app.log",
		"testdata/app.log", "testdata/app.log",
	}
	es, errs := ParseConcurrent(paths)

	if len(es) != 13*4 {
		t.Errorf("%d entradas, quiero %d", len(es), 13*4)
	}
	if len(errs) != 2*4 {
		t.Errorf("%d errores, quiero %d", len(errs), 2*4)
	}
}
*/
