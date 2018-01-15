import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PhotoListPreviewComponent } from './photo-list-preview.component';

describe('PhotoListPreviewComponent', () => {
  let component: PhotoListPreviewComponent;
  let fixture: ComponentFixture<PhotoListPreviewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PhotoListPreviewComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PhotoListPreviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
