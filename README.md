Bubble Bath
===========
Bubble Bath is a component framework for [Charm's excellent BubbleTea framework](https://github.com/charmbracelet/bubbletea) that's intended to make writing components (Bubbles) easier, particularly for fullscreen TUIs.

Get Started
-----------
Write a component that satisfies `InteractiveComponent` interface:

```go
import bubble_bath "github.com/mieubrisse/bubble-bath"

type MyApp interface {
    bubble_bath.InteractiveComponent
}
```

```go
import bubble_bath "github.com/mieubrisse/bubble-bath"

// Implementation of MyApp
type implementation interface {
    bubble_bath.InteractiveComponent
}

func New() MyApp {

}
```

Then use it in your `main.go`:

Tips
----
- For each component you create, create a public interface and a private implementation
- Keep each component in its own package (directory)

What's Inside
-------------
1. `RunBubbleBathProgram`, a wrapper over `tea.NewProgram().Run()` with sane defaults (e.g. handles resizes and quit events out of the box)
1. If you'd prefer not to use `RunBubbleBathProgram`, a `NewBubbleBathModel` function to create a `tea.Model` for use with `tea.NewProgram`
1. A `Component` interface with standardized `View`, `Resize`, `GetHeight`, and `GetWidth` functions
1. An `InteractiveComponent` interface with:
    1. A by-reference `Update(msg tea.Msg)` function, so component updating is by-reference. This sacrifices pure Redux-like state machine transitioning, but I don't need/use that right now and should make everything faster (because less by-value copying). If I need the Redux-like state machine transitioning I'll figure out a way to do it.
    1. Standardized `SetFocus` and `IsFocused` functions
1. Several out-of-the-box components conforming to `Component` that can be used to build other components:
    1. Flexbox, which allows mixed fixed-size and flexing items
    1. Text block
    1. Text input
    1. Text area
    1. Text area with Vim bindings
    1. Filterable list (which can handle nested inputs)
    1. Filterable checklist
1. Several helper methods (e.g. `GetMinInt`, `GetMaxInt`, etc.)

Why?
----
During the course of building [a decently complex TUI](https://github.com/mieubrisse/cli-journal-go) using BubbleTea, I found the vanilla BubbleTea framework useful, but difficult to work with for several reasons:

1. When I started with BubbleTea, it seemed like all my custom components should implement the `tea.Model` interface. However, I found it suitable only for the top-level model that gets slotted into `tea.NewProgram`, because I hit problems when I tried following the same pattern for subcomponents:
    1. The `tea.Model.Update` command returns `tea.Model` by-value. However, this means that you need a force-cast when calling `Update` on a subcomponent implementing the `tea.Model.Update` signature, because the subcomponent will only return `tea.Model` (not itself). For example:
       ```go
       type Parent struct{
           child Child
       }

       func (parent Parent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
           var cmd tea.Cmd
           parent.child, cmd = parent.child.Update(msg).(Child)    // <--- This cast is necessary because we're conforming to tea.Model 
           return parent, cmd
       }

       ```
       It seems that the Charm team hit the same thing, because the Bubbles in the example repository don't conform to `tea.Model` either.
    1. The by-value `Update` is also problematic when trying to create a generic component. For example, I was writing `FilterableList[T].Update`, with `T` being the element component that the list would contain. No matter how I tried, I couldn't get implementations of the `FilterableList[T]` interface to conform to the `Update(msg tea.Msg) T` function on the interface (though a better Go programmer than I may be able to).
    1. I never needed `Init()`, and my default instinct - to use it to initialize a new component's state - was wrong.
1. The concept of "focusable component" is very useful and showed up in nearly all the example Bubbles, but it's not encoded in the BubbleTea framework in any way (all the example Bubbles recreate `Focus`, `Blur`, and `Focused` by hand).
1. A resize of my terminal window should have each parent resizing their children (because the parent knows what size the children should be), but there was no out-of-the-box way for components to do this.
1. I needed `GetMinInt`, `GetMaxInt`, and `Clamp` everywhere, but BubbleTea doesn't provide this. Instead, the example Bubbles each reimplement these as private methods where needed.

Unsolved Problems
-----------------
These are problems this system doesn't yet solve but I'd like it to:

- Sizing is strictly top-down: the parent component receives a message, and it tells children what size they should be. There's no way for children components to suggest sizes back up the tree, like the web has with intrinsic vs extrinsic sizes. This would be particularly useful with the flexbox component.
- Due to everything in BubbleTea being strings, the layout of a component (width, height, padding, margin) and its styling (colors, bold, etc.) are deeply coupled. It seems like these should be decoupled - maybe by building in a DOM-like abstraction with the terminal equivalent of CSS.

Aside: as I built this, I (a backend programmer) started to deeply grok the web.
