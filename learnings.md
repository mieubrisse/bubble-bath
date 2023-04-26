- https://developer.mozilla.org/en-US/docs/Learn/CSS/Building_blocks/Sizing_items_in_CSS
    - block-level elements always span the full width they can
    - {min,max}-{width,height} u
- https://stackoverflow.com/questions/6908846/optimal-flexible-box-layout-algorithm
- https://developer.mozilla.org/en-US/docs/Learn/CSS/Building_blocks/The_box_model
    - **Block elements generate line breaks before and after themselves**
    - **Inline elements do NOT generate linebreaks before and after themselves** The next element will be on the same line if there is space
    - Inner vs outer display types
        - Inner: layout of children
        - Outer: element's participation in flow layout
- If tea gets a Way Too Long text, it will just truncate (this is good! no extra magic!)
- https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Flexible_Box_Layout
    - There are different layout modes, and "normal flow" is actually separate from "flexbox" and "grid"
- https://medium.com/@madhur.taneja/css-layouts-cced6c7a8764
    - Position & Float are deprecated (though maybe they're used for modals??)
    - Flex & Grid seem to be the future
    - We can switch between different layout algorithms/modes
- https://www.smashingmagazine.com/2016/11/css-grids-flexbox-box-alignment-new-layout-standard/
- https://www.joshwcomeau.com/css/understanding-layout-algorithms/
    - Treasure trove of a dive into how CSS layout works!!!!
    - We can opt in to different layout algos
        - `position` key sets Positioned layout
        - `float` key sets Float layout
        - otherwise, it's the parent (e.g. `display: flex`)
    - **CSS properties actually don't do anything - it's the layout algo that defines what they do!**
        - E.g. the flex algo implements `z-index`, while normal flow does not (because it's very word processor-y)
    - The reason CSS is so confusing is because each property works differently in each algo (meaning different defaults & behaviour)
    - "In Flow layout, width is a hard rule. It will take up 2000px of space, consequences be damned."
    - "In the Flexbox algorithm, width is more of a suggestion."
    - **Using `display:flex` will turn inline elements into block-level elements!!!!!**
    - Setting multiple, conflicting layout orgs will use the highest-tier one (there's some invisible tier, with Positioned seeming highest)
    - Flow layout was designed for word-processing
        - Words are stacked next to each other, forming long sentences (**inline**)
            - `<span>` and `<strong>` are inline tags
        - Sentences composed together form **blocks** (e.g. paragraphs, headings, images, etc.)
            - `<p>` and `<h1>` are block elements!
    - IMAGES ARE INLINE BY DEFAULT!
- https://www.joshwcomeau.com/css/stacking-contexts/
    - I'm not ready for it yet, but this guy seems to know a ton about CSS


Layout algo
===========
### Flex
By default, per https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Flexible_Box_Layout/Controlling_Ratios_of_Flex_Items_Along_the_Main_Ax#flex_item_sizing
1. Each flex element gets space allocated to its basis (which is MaxIntrinsicWidth if `max-content`)
    - VERY useful, this link!!!!
1. When overbudgeted, elements with larger bases get more stuff removed than elements with smaller bases


1. Ask child what its intrinsic width is (max desired width)
    1. Child reports clamp(content_width, min_width, max_width)
        - Note that hard-setting the width of an item will set desired_width == min_width == max_width
1. Give child the 
1. Give child 

1. Viewport truncates to actual terminal width/height
