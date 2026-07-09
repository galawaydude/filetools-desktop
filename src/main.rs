use gpui::*;
use gpui_component::{button::*, *};

/// Temporary scaffold view — confirms GPUI + gpui-component render end to end.
struct App;

impl Render for App {
    fn render(&mut self, _: &mut Window, cx: &mut Context<Self>) -> impl IntoElement {
        div()
            .v_flex()
            .gap_4()
            .size_full()
            .items_center()
            .justify_center()
            .bg(cx.theme().background)
            .child("File Tools")
            .child(
                Button::new("go")
                    .primary()
                    .label("Let's Go!")
                    .on_click(|_, _, _| println!("Clicked!")),
            )
    }
}

fn main() {
    gpui_platform::application().run(move |cx| {
        gpui_component::init(cx);

        cx.spawn(async move |cx| {
            cx.open_window(WindowOptions::default(), |window, cx| {
                let view = cx.new(|_| App);
                cx.new(|cx| Root::new(view, window, cx).bg(cx.theme().background))
            })
            .expect("failed to open window");
        })
        .detach();
    });
}
