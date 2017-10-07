import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpModule } from '@angular/http';

import { AppRoutesModule } from './app-routes.module';

import { CollectionService } from './services/collection.service';
import { PathService  } from './services/path.service';
import { PhotoService } from './services/photo.service';
import { RenditionConfigurationService } from './services/rendition-configuration.service';

import { AppComponent } from './app.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { CollectionComponent } from "./collections/collection/collection.component";
import { CollectionDashboardComponent } from "./collections/collection-dashboard/collection-dashboard.component";

// import { DashboardComponent } from './collections/dashboard/dashboard.component';
// import { ListComponent } from './collections/list/list.component';
// import { DetailsComponent } from './collections/details/details.component';
// import { PhotoDetailsComponent } from "./photos/photo-details/photo-details.component";
// import { CollectionsComponent } from "./collections/collections.component";
// import { FrameComponent } from "./collections/frame/frame.component";

@NgModule({
  declarations: [
    AppComponent,
    DashboardComponent,
    CollectionComponent,
    CollectionDashboardComponent
  ],
  imports: [
    BrowserModule,
    HttpModule,
    AppRoutesModule
  ],
  providers: [
    CollectionService,
    PhotoService,
    PathService,
    RenditionConfigurationService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
