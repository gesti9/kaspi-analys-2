#!/bin/bash

# Функция для замены содержимого файла на нули
replace_with_zero_in_file() {
    local file=$1
    local length=$(wc -c < "$file")  # Определение длины содержимого файла
    printf '0%.0s' $(seq 1 $length) > "$file"  # Запись нулей в файл
}

# Проход по всем .txt файлам в текущей директории и её поддиректориях
process_current_directory() {
    find . -type f -name "*.txt" | while read file; do
        replace_with_zero_in_file "$file"
    done
}

process_current_directory

