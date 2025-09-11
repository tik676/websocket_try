import tkinter as tk
from PIL import Image, ImageTk

win = tk.Tk()  

win.title("nazvanie")

win.geometry("500x500")

win.resizable(True,False) # Запрещает менять размер

win.maxsize(700,800)


icon = ImageTk.PhotoImage(Image.open("icon.png"))
win.iconphoto(False, icon)
print("Иконка успешно загружена!")

label1 = tk.Label(win,text="LABEL",fg="blue",bg="green")
label1.pack()


win.mainloop()  # Главный цикл
