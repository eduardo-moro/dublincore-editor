<img width="1015" height="556" alt="image" src="https://github.com/user-attachments/assets/c9ea5d39-abe0-487b-9458-257588ef1427" />

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-square&logo=go)
![CLI](https://img.shields.io/badge/Interface-CLI-brightgreen)
![License](https://img.shields.io/badge/Licença-MIT-blue)

Uma ferramenta em Go para editar metadados Dublin Core em arquivos DOCX com uma interface TUI (Text User Interface) amigável usando BubbleTea.

## ✨ Características

- **Metadados ATS**: Foco em metadados para Applicant Tracking Systems
- **Backup Automático**: Cria backup automático antes de editar
- **Suporte a DOCX**: Compatível com arquivos Microsoft Word Originais
- **Campos Essenciais**: Edição dos 5 campos [mais importantes para currículos](https://www.youtube.com/watch?v=fQ7GMBIDric)

## 🚀 Instalação

### Pré-requisitos

- [Go](https://golang.org/dl/) 1.21 ou superior
- Arquivos DOCX do Microsoft Word (não Google Docs)

### Instalação Rápida

```bash
# Clone o repositório
git clone https://github.com/eduardo-moro/metadata-editor.git
cd metadata-editor

# Instale as dependências
go mod download

# Construa o executável
go build -o dcedit.exe main.go

# Adicione ao PATH (Windows)
move dcedit.exe %USERPROFILE%\bin\
setx PATH "%PATH%;%USERPROFILE%\bin"
```

## 📋 Uso

### Editar Metadados (Interface Visual)
```bash
dcedit "C:\caminho\para\seu\curriculo.docx"
```

### Visualizar Metadados Atuais
```bash
dcedit view --file "C:\caminho\para\seu\curriculo.docx"
```

### Debug do Arquivo (Para Desenvolvedores)
```bash
dcedit debug --file "C:\caminho\para\seu\curriculo.docx"
```

### Especificar Arquivo de Saída
```bash
dcedit edit --output "C:\caminho\para\curriculo_atualizado.docx" "C:\caminho\para\seu\curriculo.docx"
```

## 🎯 Campos de Metadados Suportados

### 1. **DC: Title** (Título)
- Exemplo: "Analista Backend Pleno"
- O cargo principal ou título do currículo

### 2. **DC: Creator** (Criador)
- Exemplo: "Eduardo Moro"
- Nome do candidato (separado por vírgulas se múltiplos)

### 3. **CP: Keywords** (Palavras-chave)
- Exemplo: "Go, PHP, AWS, Docker, Kubernetes"
- Habilidades e tecnologias (separadas por vírgulas)

### 4. **CP: Description** (Descrição)
- Exemplo: "Analista backend com 6 anos de experiência em tecnologia..."
- Resumo profissional ou objetivo

### 5. **CP: Category** (Categoria)
- Valor fixo: "curriculo"
- Otimizado para sistemas ATS

## 🛠️ Para Desenvolvedores

### Estrutura do Projeto
```
metadata-editor/
├── main.go                 # Ponto de entrada principal
├── go.mod                 # Dependências do Go
├── ui/
│   └── editor.go          # Interface BubbleTea TUI
├── docx/
│   └── docx.go           # Manipulação de arquivos DOCX
├── dublincore/
│   └── dublincore.go     # Modelos de metadados Dublin Core
└── cmd/
    └── editor/
        └── editor.go     # Comandos CLI
```

### Dependências Principais
- [BubbleTea](https://github.com/charmbracelet/bubbletea): TUI framework
- [CLI](https://github.com/urfave/cli): Framework de linha de comando
- [unioffice](https://github.com/unidoc/unioffice): Manipulação de documentos Office

### Build e Desenvolvimento
```bash
# Build para desenvolvimento
go build -o dcedit.exe main.go

# Executar testes
go test ./...

# Limpar build
go clean
```

## ❗ Limitações e Requisitos

### ✅ Suportado
- Arquivos DOCX do Microsoft Word
- Metadados Dublin Core e Core Properties
- Encoding UTF-8
- Sistemas Windows, Linux e macOS

### ❌ Não Suportado
- Arquivos do Google Docs exportados como DOCX
- Documentos protegidos por senha
- Formatos antigos (.doc)
- Metadados personalizados não padrão

## 🔧 Solução de Problemas

### Erro: "zip: not a valid zip file"
- **Causa**: Arquivo corrompido ou não é um DOCX válido
- **Solução**: Abra e salve o arquivo no Microsoft Word

### Erro: "Google Docs export detected"
- **Causa**: Arquivo exportado do Google Docs
- **Solução**: Salve o arquivo usando "Salvar como" no Microsoft Word

### Metadados não aparecem após edição
- **Causa**: Problema de parsing do XML
- **Solução**: Use `dcedit debug --file arquivo.docx` para diagnosticar

## 📝 Exemplos de Uso

### Exemplo 1: Edição Completa
```bash
# Edita um currículo com interface visual
dcedit "C:\Curriculos\meu_curriculo.docx"
```

### Exemplo 2: Batch Processing
```bash
# Script para processar múltiplos arquivos
@echo off
for %%i in (C:\Curriculos\*.docx) do (
    echo Processando: %%~ni
    dcedit "%%i" --output "C:\Curriculos\Processados\%%~ni.docx"
)
```

## 🤝 Contribuindo

1. Faça um Fork do projeto
2. Crie uma Branch para sua Feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas Mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a Branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🆘 Suporte

Se você encontrar problemas:

1. Verifique se o arquivo é um DOCX válido do Microsoft Word
2. Execute `dcedit debug --file arquivo.docx` para diagnosticar
3. Abra uma issue no GitHub com:
   - Arquivo de exemplo (se possível)
   - Comando executado
   - Saída do debug

---



⭐️ Se este projeto foi útil, deixe uma estrela no GitHub!


<div align="right">
    <span title="Com amor, não com lovable">Feito com ❤️</span>
</div>
