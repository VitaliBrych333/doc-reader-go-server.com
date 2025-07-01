package documents

type Info struct {
    Comments      string      `json:"comments"`
    Author        string      `json:"author"`
    Date_Created  string      `json:"dateCreated"`
    Date_Modified string      `json:"dateModified"`
}

type RespDocument struct {
    Id            int         `json:"id"`
    User_Id       string      `json:"userId"`
    Document_Id   string      `json:"documentId"`
    Document_Name string      `json:"name"`
    File          []byte      `json:"file"`
    Info          Info        `json:"info"`
} 

type ReqDocument struct {
    Id            int         `form:"id"`
    User_Id       string      `form:"userId"`
    Document_Id   string      `form:"documentId"`
    Document_Name string      `form:"name"`
    File          []byte      `form:"file"`
    Info          Info        `form:"info"`

}
type PageInfo struct {
    Page_Id                string      `json:"pageId"`
    Num_Page               int         `json:"numPage"`
    Rotate                 int         `json:"rotate"`
    Original_Document_Id   string      `json:"originalDocumentId"`
    Original_Num_Page      int         `json:"originalNumPage"`
}

type Page struct {
    Id            string      `json:"id"`
    Page          PageInfo    `json:"page"`
}

type Document struct {
    Id            string      `json:"id"`
    Name          string      `json:"name"`
    Pages         []PageInfo  `json:"pages"`
}

type DeletePageAction struct {
    Id            string       `json:"id"`
    Page          PageInfo     `json:"page"`
}

type NewDocumentAction struct {
    Doc             Document    `json:"doc"`
    Position_Index  int         `json:"positionIndex"`
}

type NewPageAction struct {
    Position_Index  int         `json:"positionIndex"`
    Page            PageInfo    `json:"page"`
}

type RenameAction struct {
    Id              string      `json:"id"`
    Name            string      `json:"name"`
}

type EditActionValue struct {
    Id               string            `json:"id,omitempty"`
    Name             string            `json:"name,omitempty"`
    Position_Index   int               `json:"positionIndex,omitempty"`
    Doc              Document          `json:"doc,omitempty"`
    Page             PageInfo          `json:"page,omitempty"`
    Info             Info              `json:"info,omitempty"`
}

type EditAction struct {
    Type             string            `json:"type"`
    Value            EditActionValue   `json:"value"`
}
type FileDocument struct {
    Id               string            `json:"id"`
    Document_Name    string            `json:"name"`
    Info             Info              `json:"info"`
    File             []byte            `json:"file"`
    Pages            []PageInfo        `json:"pages"`
}

type ReqSaveDocuments struct {
    User_Id          string            `json:"userId"`
    New_Documents    []FileDocument    `json:"newDocuments"`
    Edit_Actions     []EditAction      `json:"editActions"`
    Rotate           []Page            `json:"rotate"`
}