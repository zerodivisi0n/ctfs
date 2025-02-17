# Exploiting Metadata Processing with ExifTool

## 📜 Scenario
You have access to a server with two endpoints that process documents and modify their metadata using **ExifTool**:

1. **PDF Metadata Injection**
   - **Endpoint:** `POST /pdf?title=sometitle`
   - **Behavior:** Takes a PDF file, updates its **EXIF metadata title** based on the `title` query parameter, and returns the modified PDF.

2. **DOCX to PDF Conversion with Metadata Handling**
   - **Endpoint:** `POST /docx`
   - **Behavior:** Accepts a **DOCX file**, converts it to **PDF**, extracts the title from the **DOCX document properties**, and embeds it into the **PDF's EXIF metadata**.

## 🎯 Objective
Your task is to **exploit vulnerabilities** in this system to achieve the following:

1. **File Read** – Extract sensitive files from the server.
2. **File Write** – Modify or create arbitrary files.
3. **Remote Code Execution (RCE)** – Execute commands on the server.

## 🛠️ Setup Instructions
Start the vulnerable server using one of the following methods:

- **Using Docker:**
    ```sh
    docker run -it --rm -p 3000:3000 ghcr.io/zerodivisi0n/exiftoolpwn
    ```
- **Building and Running Locally:**
    ```sh
    make run
    ```

## 🔥 Exploitation
Use `curl` to interact with the server:

- **Modify metadata in a PDF:**
    ```sh
    curl -X POST "localhost:3000/pdf?title=Hacked" --data-binary @document.pdf -o result.pdf
    ```
- **Send a DOCX file for conversion:**
    ```sh
    curl -X POST localhost:3000/docx --data-binary @document.docx -o result.pdf
    ```

Can you **exploit** this system and **take control**? Good luck! 🚀