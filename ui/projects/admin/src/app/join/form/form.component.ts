import { EventEmitter, Component, OnInit, Input, Output } from '@angular/core';

import { JoinRequest, Invitation } from '../join.component';

@Component({
  selector: 'app-form',
  templateUrl: './form.component.html',
  styleUrls: ['./form.component.css'],
})
export class FormComponent implements OnInit {
  @Input() invitation: Invitation;

  @Input() disabled = false;

  @Output() submitted = new EventEmitter<JoinRequest>();

  joinRequest: JoinRequest = new JoinRequest();

  constructor() {}

  ngOnInit(): void {
    this.joinRequest.email = this.invitation.email;
    this.joinRequest.token = this.invitation.token;
  }

  onSubmit() {
    this.submitted.emit(this.joinRequest);
  }
}
