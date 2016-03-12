package main

const CSS_CODE = `
  GtkWindow {
    border: 1px solid rgba(0,0,0,0.2);
  }
  GtkEntry {
    background-image: none;
    border: none;
    box-shadow: none;
    font-size: 12pt;
  }
  GtkTreeView {
    background-color: transparent;
    font-size: 8pt;
  }
  GtkTreeView:selected {
    background-color: rgba(0,0,0,0.1);
  }
  GtkTreeView.cell {
    padding: 6px 3px;
  }
`
