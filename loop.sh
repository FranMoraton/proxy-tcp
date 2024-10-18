#!/bin/bash

# Ejecutar 100 veces `make curl` de manera concurrente
for i in {1..100}; do
    echo "Llamada $i"
    make curl &   # Ejecuta make curl en background
done

# Esperar a que todas las llamadas terminen
wait

echo "Todas las llamadas curl han terminado."
