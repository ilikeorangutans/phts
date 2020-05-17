import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { NotFoundComponent } from './not-found/not-found.component';
import { LandingPageComponent } from './landing-page/landing-page.component';
import { ShareViewerComponent } from './share-viewer/share-viewer.component';

@NgModule({
  declarations: [
    AppComponent,
    NotFoundComponent,
    LandingPageComponent,
    ShareViewerComponent,
  ],
  imports: [BrowserModule, AppRoutingModule, HttpClientModule],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
