import { LandingComponent } from './landing/landing.component';
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';


import { PublicRoutingModule } from './public-routing.module';
import { AppComponent } from './app/app.component';
import { ShareViewerComponent } from './share-viewer/share-viewer.component';
import { PathService } from './services/path.service';
import { ShareService } from './services/share.service';

@NgModule({
  imports: [
    CommonModule,
    PublicRoutingModule
  ],
  declarations: [
    LandingComponent,
    AppComponent,
    ShareViewerComponent
  ],
  providers: [
    PathService,
    ShareService
  ]
})
export class PublicModule { }
