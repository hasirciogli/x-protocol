import os

# Bu dosya ve klasörler hariç tutulacak
EXCLUDED_FOLDERS = {".git", ".idea", "node_modules", "__pycache__"}

def generate_tree(path, indent="", output_file=None):
    # Eğer dizinse
    if os.path.isdir(path):
        folder_name = os.path.basename(path)
        
        # Eğer bu klasör hariç tutulmuşsa, işlemi geç
        if folder_name in EXCLUDED_FOLDERS:
            return
        
        # Çıktıya klasörün adını yaz
        if output_file:
            output_file.write(f"{indent}[DIR] {folder_name}\n")
        else:
            print(f"{indent}[DIR] {folder_name}")
        
        # Dizin içindeki her şeyi dolaş
        for item in os.listdir(path):
            generate_tree(os.path.join(path, item), indent + "    ", output_file)
    else:
        # Eğer dosya ise, dosyanın adını yaz
        if output_file:
            output_file.write(f"{indent}{os.path.basename(path)}\n")
        else:
            print(f"{indent}{os.path.basename(path)}")

# Projenin kök dizini
project_dir = "."  # Bu, geçerli dizini temsil eder, istersen değiştirebilirsin

# Çıktıyı bir txt dosyasına yazalım
with open("project_tree.txt", "w") as f:
    generate_tree(project_dir, output_file=f)

print("Proje yapısı project_tree.txt dosyasına kaydedildi.")
