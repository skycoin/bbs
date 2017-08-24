import { NgModule, ModuleWithProviders } from '@angular/core';
import { PopupWindow } from './popup-window';
import { Popup } from './popup';
import { PopupStack } from './popup-stack';
@NgModule({
  imports: [],
  exports: [],
  declarations: [PopupWindow],
  entryComponents: [PopupWindow],
  providers: [],
})
export class PopupModule {
  static forRoot(): ModuleWithProviders { return { ngModule: PopupModule, providers: [Popup, PopupStack] } };
}
