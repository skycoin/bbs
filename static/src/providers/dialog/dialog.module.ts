import { NgModule, ModuleWithProviders } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DialogComponent } from './dialog.component';
import { Dialog } from './dialog.service';
import { DialogOverlayComponent } from './dialog.overlay.component';

@NgModule({
  imports: [CommonModule],
  exports: [],
  declarations: [DialogComponent, DialogOverlayComponent],
  entryComponents: [DialogComponent, DialogOverlayComponent],
  providers: [],
})
export class DialogModule {
  static forRoot(): ModuleWithProviders { return { ngModule: DialogModule, providers: [Dialog] } };
}
