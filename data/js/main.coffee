FileListRequestProtocol = "fsrp"
url = "ws://"+window.location.host+"/ws"
webSocket = new WebSocket(url)
fileList = null
paths = new Array()

webSocket.onopen = () ->
    webSocket.send(FileListRequestProtocol)
    return

webSocket.onmessage = (event) ->
    message = JSON.parse(event.data)
    for messageType, messageContent of message
        switch messageType
            when FileListRequestProtocol
                fileList = message[FileListRequestProtocol]
                updateTable()
    return

updateTable = () ->
    list = fileList
    for dir in paths
        if list["Children"][dir]?
            list = list["Children"][dir]
        else
            list = fileList
            paths = new Array()
            break
    table = document.getElementById("fileList")
    for i in [1...table.rows.length]
        table.deleteRow(1)

    i = 1
    for name, file of list["Children"]
        row = table.insertRow(i)
        i++
        addFileToRow(row, file)
    $("tbody > tr").off().click(onClickRow)
    return

addFileToRow = (row, file) ->
    name = row.insertCell(0)
    isDir = row.insertCell(1)
    atLocal = row.insertCell(2)
    size = row.insertCell(3)

    name.innerHTML = file.Name
    isDir.innerHTML = file.IsDir
    atLocal.innerHTML = file.AtLocal
    size.innerHTML = file.Size
    return

onClickRow = () ->
    name = this.getElementsByTagName("td")[0].innerHTML
    isDir = this.getElementsByTagName("td")[1].innerHTML
    if name == ".." and paths.length > 0
        paths.pop()
    else if isDir == "true"
        paths.push(name)
    updateTable()
    return

