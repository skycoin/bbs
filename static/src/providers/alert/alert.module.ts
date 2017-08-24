import { NgModule, ModuleWithProviders } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AlertWindowComponent } from './alert-window';
import { AlertComponent } from './alert';
import { Alert } from './alert.service';
// import { PopupStack } from './popup-stack';
@NgModule({
  imports: [CommonModule],
  exports: [],
  declarations: [AlertWindowComponent, AlertComponent],
  entryComponents: [AlertWindowComponent, AlertComponent],
  providers: [],
})
export class AlertModule {
  static forRoot(): ModuleWithProviders { return { ngModule: AlertModule, providers: [Alert] } };
}
