import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { RenditionConfigurationsComponent } from './rendition-configurations.component';

describe('RenditionConfigurationsComponent', () => {
  let component: RenditionConfigurationsComponent;
  let fixture: ComponentFixture<RenditionConfigurationsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ RenditionConfigurationsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(RenditionConfigurationsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
