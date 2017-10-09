import { UploadQueueService } from './services/upload-queue.service';
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpModule } from '@angular/http';
import { FormsModule } from '@angular/forms';
import { AppRoutesModule } from './routes/app-routes.module';

import { CollectionService } from './services/collection.service';
import { PathService  } from './services/path.service';
import { PhotoService } from './services/photo.service';
import { RenditionConfigurationService } from './services/rendition-configuration.service';

import { AppComponent } from './app.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { CollectionComponent } from './collections/collection/collection.component';
import { CollectionDashboardComponent } from './collections/collection-dashboard/collection-dashboard.component';
import { PhotoUploadComponent } from './photos/photo-upload/photo-upload.component';
import { CollectionCreateComponent } from './collections/collection-create/collection-create.component';
import { PhotoDetailsComponent } from './photos/photo-details/photo-details.component';
import { CollectionBrowserComponent } from './collections/collection-browser/collection-browser.component';
import { UploadQueueComponent } from './photos/upload-queue/upload-queue.component';
import { ConfigurationListComponent } from './renditions/configuration-list/configuration-list.component';

@NgModule({
  declarations: [
    AppComponent,
    DashboardComponent,
    CollectionComponent,
    CollectionDashboardComponent,
    CollectionCreateComponent,
    PhotoUploadComponent,
    PhotoDetailsComponent,
    CollectionBrowserComponent,
    UploadQueueComponent,
    ConfigurationListComponent
  ],
  imports: [
    BrowserModule,
    HttpModule,
    AppRoutesModule,
    FormsModule
  ],
  providers: [
    CollectionService,
    PhotoService,
    PathService,
    RenditionConfigurationService,
    UploadQueueService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
