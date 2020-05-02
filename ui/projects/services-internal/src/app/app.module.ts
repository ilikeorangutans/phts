import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { BasePathService } from 'projects/shared/src/app/services/base-path.service';
import { PathService } from './services/path.service';
import { AuthService } from './services/auth.service';

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule
  ],
  providers: [
    AuthService,
    BasePathService,
    PathService,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
