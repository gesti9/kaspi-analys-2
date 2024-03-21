from selenium import webdriver
from selenium.webdriver.common.by import By

# Укажите путь к вашему драйверу
driver_path = "C:\chromedriver_win32\chromedriver.exe"

# Создайте экземпляр драйвера
driver = webdriver.Chrome(executable_path=driver_path)

# Откройте страницу
url = "https://kaspi.kz/shop/p/elektrochainik-bereke-br-810-seryi-109981423/?c=710000000"
driver.get(url)

# Найдите элемент с классом item__price-once
element = driver.find_element(By.CLASS_NAME, "item__price-once")

# Получите текст элемента
element_text = element.text

# Выведите текст элемента
print(element_text)

# Закройте браузер
driver.quit()
