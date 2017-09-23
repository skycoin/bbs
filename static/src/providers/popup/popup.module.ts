import { NgModule, ModuleWithProviders } from '@angular/core';
import { CommonModule } from '@angular/common'
import { PopupWindow } from './popup-window';
import { PopupBackdrop } from './popup.backdrop';
import { Popup } from './popup';
import { PopupStack, ActivePop } from './popup-stack';
@NgModule({
  imports: [],
  exports: [],
  declarations: [PopupWindow, PopupBackdrop],
  entryComponents: [PopupWindow, PopupBackdrop],
  providers: [],
})
export class PopupModule {
  static forRoot(): ModuleWithProviders { return { ngModule: PopupModule, providers: [Popup, PopupStack, ActivePop] } };
}
