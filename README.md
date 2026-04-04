# tui-images

TUI (Terminal User Interface) para comprimir imágenes de forma rápida y sencilla desde la terminal.

## Características

- Selección interactiva de archivos o carpetas
- Compresión de imágenes JPEG y PNG
- Selección automática del formato más pequeño (compara JPEG vs PNG)
- Calidad configurable (1-100)
- Barra de progreso en tiempo real
- Resumen detallado: tamaño original, comprimido y porcentaje ahorrado
- Soporte recursivo para carpetas con subdirectorios
- Instalación global — ejecutable desde cualquier directorio

## Requisitos

- **Go 1.21+** instalado en el sistema
- Terminal compatible (soporta secuencias de escape ANSI)

## Instalación

### Desde el código fuente

```bash
# Clonar o descargar el repositorio
cd tui-images

# Instalar dependencias
go mod tidy

# Compilar
go build -o tui-images ./cmd/main.go

# Instalar globalmente (opcional)
sudo cp tui-images /usr/local/bin/
```

### Con `go install`

```bash
go install github.com/fran-codigo/tui-images/cmd@latest
```

## Uso

### TUI interactiva

```bash
tui-images              # Calidad por defecto 75%
tui-images -q 50        # Calidad 50%
```

### Flujo de uso

1. **Seleccionar modo**: `F` para archivo individual, `D` para carpeta
2. **Navegar**: Flechas `↑` `↓` para moverse, `Enter` para seleccionar
3. **Calidad**: Escribir número del 1 al 100, `Enter` para comprimir
4. **Resultado**: Se muestra resumen con tamaños y porcentaje ahorrado

### Atajos de teclado

| Tecla | Acción |
|-------|--------|
| `F` / `D` | Modo archivo / directorio |
| `Tab` | Cambiar entre modos |
| `↑` `↓` | Navegar archivos |
| `Enter` | Confirmar selección |
| `Esc` | Volver atrás |
| `Q` / `Ctrl+C` | Salir |

### Flags

| Flag | Descripción | Default |
|------|-------------|---------|
| `-q` | Calidad de compresión (1-100) | 75 |
| `-v` | Mostrar versión | — |

## Formatos soportados

| Entrada | Salida | Notas |
|---------|--------|-------|
| `.jpg`, `.jpeg` | `.jpg` | Compresión JPEG con calidad configurable |
| `.png` | `.jpg` o `.png` | Se elige automáticamente el formato más pequeño |

> **Nota**: Para PNG se prueba tanto JPEG (con la calidad dada) como PNG (máxima compresión) y se guarda el que resulte en menor tamaño.

## Estructura del proyecto

```
tui-images/
├── cmd/
│   └── main.go              # Entry point, parse de flags
├── internal/
│   ├── compressor/
│   │   └── compressor.go    # Lógica de compresión
│   └── tui/
│       ├── model.go         # Modelo Bubbletea
│       ├── update.go        # Manejo de eventos
│       └── view.go          # Renderizado UI con Lipgloss
├── go.mod
└── go.sum
```

## Tecnologías

| Componente | Librería |
|------------|----------|
| TUI Framework | [Bubbletea](https://github.com/charmbracelet/bubbletea) |
| Componentes UI | [Bubbles](https://github.com/charmbracelet/bubbles) (filepicker, progress) |
| Estilos | [Lipgloss](https://github.com/charmbracelet/lipgloss) |
| Compresión | `image/jpeg`, `image/png` (stdlib de Go) |

## Seguridad

- No modifica ni borra archivos originales
- Crea copias comprimidas en una carpeta nueva (`compressed/` o `nombre_compressed/`)
- Rechaza symlinks para evitar acceso fuera del directorio seleccionado
- Valida que cada archivo esté dentro del directorio seleccionado (previene path traversal)
- Límite de 100MB por archivo (previene OOM)
- Permisos restrictivos en archivos creados (`0640`)
- Sin dependencias externas para procesamiento de imágenes (usa solo stdlib)

## Licencia

MIT
