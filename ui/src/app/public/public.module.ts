import { LandingComponent } from './landing/landing.component';
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { PublicRoutingModule } from './public-routing.module';
import { AppComponent } from './app/app.component';
import { ShareViewerComponent } from './share-viewer/share-viewer.component';

@NgModule({
  imports: [
    CommonModule,
    PublicRoutingModule
  ],
  declarations: [
    LandingComponent,
    AppComponent,
    ShareViewerComponent
  ]
})
export class PublicModule { }
