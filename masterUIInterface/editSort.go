package masterUIInterface

import (
	"fmt"
  "log"
  "github.com/jroimartin/gocui"
  "github.com/kkellner/cloudfoundry-top-plugin/util"
)

const MAX_SORT_COLUMNS = 4

type EditSortView struct {
  masterUI MasterUIInterface
	name string
  width int
  height int
  listWidget *ListWidget

  sortPosition int
  sortColumns []*sortColumn
}

func NewEditSortView(masterUI MasterUIInterface, name string, listWidget *ListWidget) *EditSortView {
	w := &EditSortView{masterUI: masterUI, name: name, listWidget: listWidget}
  w.width = 55
  w.height = 10

  w.sortColumns = make([]*sortColumn,MAX_SORT_COLUMNS)

  for i, sc := range listWidget.sortColumns {
    w.sortColumns[i] = sc
  }
  return w
}

func (w *EditSortView) Layout(g *gocui.Gui) error {
  maxX, maxY := g.Size()
	v, err := g.SetView(w.name, maxX/2-(w.width/2), maxY/2-(w.height/2), maxX/2+(w.width/2), maxY/2+(w.height/2))
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
    v.Title = "Edit Sort"
    v.Frame = true
    fmt.Fprintln(v, "...")
    if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.closeView); err != nil {
      return err
    }
    if err := g.SetKeybinding(w.name, gocui.KeyArrowRight, gocui.ModNone, w.keyArrowRightAction); err != nil {
      return err
    }
    if err := g.SetKeybinding(w.name, gocui.KeyArrowLeft, gocui.ModNone, w.keyArrowLeftAction); err != nil {
      return err
    }
    if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.keyArrowDownAction); err != nil {
      return err
    }
    if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.keyArrowUpAction); err != nil {
      return err
    }
    if err := g.SetKeybinding(w.name, gocui.KeySpace, gocui.ModNone, w.keySpaceAction); err != nil {
      return err
    }
    if err := g.SetKeybinding(w.name, gocui.KeyDelete, gocui.ModNone, w.keyDeleteAction); err != nil {
      return err
    }
    if err := g.SetKeybinding(w.name, gocui.KeyBackspace, gocui.ModNone, w.keyDeleteAction); err != nil {
      return err
    }
    if err := g.SetKeybinding(w.name, gocui.KeyBackspace2, gocui.ModNone, w.keyDeleteAction); err != nil {
      return err
    }

    if err := w.masterUI.SetCurrentViewOnTop(g, w.name); err != nil {
      log.Panicln(err)
    }
    w.RefreshDisplay(g)
	}
	return nil
}

func (w *EditSortView) RefreshDisplay(g *gocui.Gui) error {
  v, err := g.View(w.name)
  if err != nil {
    return err
  }
  v.Clear()
  fmt.Fprintln(v, " Right Arrow or Left Arrow to select sort column,")
  fmt.Fprintln(v, " press SPACE to select column for sorting.")
  fmt.Fprintln(v, " Press ENTER to apply sort")
  fmt.Fprintln(v, "")

  for i, sc := range w.sortColumns {
    if w.sortPosition == i {
      fmt.Fprintf(v, util.REVERSE_WHITE)
    }
    displayName := "-none-"
    if sc != nil {
      sortDirection := "(ascending)"
      if sc.reverseSort {
        sortDirection = "(descending)"
      }
      columnLabel := w.listWidget.columnMap[sc.id].label
      displayName = fmt.Sprintf("%-13v %v", columnLabel, sortDirection)
    }
    fmt.Fprintf(v, " Sort #%v: %v\n", i+1, displayName)
    if w.sortPosition == i {
      fmt.Fprintf(v, util.CLEAR)
    }
  }
  return w.listWidget.RefreshDisplay(g)
}


func (w *EditSortView) keyArrowRightAction(g *gocui.Gui, v *gocui.View) error {
  columnId := w.listWidget.editSortColumnId
  columns := w.listWidget.columns
  columnsLen := len(columns)
  for i, col := range columns {
    if col.id == columnId && i+1 < columnsLen {
      columnId = columns[i+1].id
      break
    }
  }
  //writeFooter(g, fmt.Sprintf("\r columnId: %v", columnId) )
  w.listWidget.editSortColumnId = columnId
  w.RefreshDisplay(g)
  return nil
}

func (w *EditSortView) keyArrowLeftAction(g *gocui.Gui, v *gocui.View) error {
  columnId := w.listWidget.editSortColumnId
  columns := w.listWidget.columns
  for i, col := range columns {
    if col.id == columnId && i > 0 {
      columnId = columns[i-1].id
      break
    }
  }
  //writeFooter(g, fmt.Sprintf("\r columnId: %v", columnId) )
  w.listWidget.editSortColumnId = columnId
  return w.RefreshDisplay(g)
}

func (w *EditSortView) keyArrowDownAction(g *gocui.Gui, v *gocui.View) error {
  if w.sortPosition+1 < MAX_SORT_COLUMNS {
    w.sortPosition++
  }
  return w.RefreshDisplay(g)
}

func (w *EditSortView) keyArrowUpAction(g *gocui.Gui, v *gocui.View) error {
  if w.sortPosition >0 {
    w.sortPosition--
  }
  return w.RefreshDisplay(g)
}

func (w *EditSortView) keyDeleteAction(g *gocui.Gui, v *gocui.View) error {
  writeFooter(g, fmt.Sprintf("\r DELETE") )
  w.sortColumns[w.sortPosition] = nil
  return w.RefreshDisplay(g)
}

func (w *EditSortView) keySpaceAction(g *gocui.Gui, v *gocui.View) error {

  sc := w.sortColumns[w.sortPosition]
  columnId := w.listWidget.editSortColumnId
  if sc == nil {
    sc = &sortColumn{
      id: columnId,
      reverseSort: w.listWidget.columnMap[columnId].defaultReverseSort,
    }
    w.sortColumns[w.sortPosition] = sc
  } else {
    if sc.id == columnId {
      sc.reverseSort = !sc.reverseSort
    } else {
      sc.id = columnId
      sc.reverseSort = w.listWidget.columnMap[columnId].defaultReverseSort
    }

  }

  writeFooter(g, fmt.Sprintf("\r sc: %+v", sc) )
  w.applySort(g)

  return nil
}

func (w *EditSortView) applySort(g *gocui.Gui) {

  useSortColumns := make([]*sortColumn,0)
  for _, sc := range w.sortColumns {
    if sc != nil {
      useSortColumns = append(useSortColumns, sc)
    }
  }
  w.listWidget.sortColumns = useSortColumns
  w.listWidget.SortData()
  w.RefreshDisplay(g)
}


func (w *EditSortView) closeView(g *gocui.Gui, v *gocui.View) error {
  w.listWidget.enableSortEdit(false)
  if err := w.masterUI.CloseView(w, w.name); err != nil {
    return err
  }
  w.applySort(g)
  //w.listWidget.RefreshDisplay(g)
	return nil
}