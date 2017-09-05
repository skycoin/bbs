import { NgModule, ModuleWithProviders } from '@angular/core';
import { LoadingComponent } from './loading';
import { LoadingService } from './loading.service';

@NgModule({
  declarations: [LoadingComponent],
  entryComponents: [LoadingComponent],
})
export class LoadingModule {
  static forRoot(): ModuleWithProviders {
    return { ngModule: LoadingModule, providers: [LoadingService] }
  };
}
