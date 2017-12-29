import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { CollectionRoutingModule } from './collection-routing.module';
import { DashboardComponent } from './dashboard/dashboard.component';
import { BrowserComponent } from './browser/browser.component';
import { AppComponent } from './app/app.component';
import { FormComponent } from './form/form.component';
import { FormsModule } from '@angular/forms';
import { LandingComponent } from './landing/landing.component';
import { RenditionConfigurationsComponent } from './rendition-configurations/rendition-configurations.component';
import { PhotoUploadComponent } from './photo-upload/photo-upload.component';
import { UploadQueueComponent } from './upload-queue/upload-queue.component';
import { PhotoStreamComponent } from './photo-stream/photo-stream.component';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    CollectionRoutingModule
  ],
  declarations: [
    DashboardComponent,
    BrowserComponent,
    AppComponent,
    FormComponent,
    LandingComponent,
    RenditionConfigurationsComponent,
    PhotoUploadComponent,
    UploadQueueComponent,
    PhotoStreamComponent
  ]
})
export class CollectionModule { }
