from flask import Flask, render_template, request, redirect, url_for
from flask_sqlalchemy import SQLAlchemy

app = Flask(__name__)

# Set up local SQLite database (UNCOMMENTED NOW)
app.config['SQLALCHEMY_DATABASE_URI'] = 'sqlite:///thrift.db'
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False
db = SQLAlchemy(app)

# 1. Database Model (Stores Member Info)
class Member(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(100), nullable=False)
    monthly_amount = db.Column(db.Float, nullable=False)
    queue_position = db.Column(db.Integer, nullable=False)  # Stores queue order (1, 2, 3...)
    has_received = db.Column(db.Boolean, default=False)

# Create the database tables automatically
with app.app_context():
    db.create_all()

# 2. Home Page Route (Show members & queue status)
@app.route('/')
def home():
    # Order members by their queue position (1, 2, 3...)
    members = Member.query.order_by(Member.queue_position.asc()).all()
    eligible_members = Member.query.filter_by(has_received=False).all()
    
    return render_template('index.html', members=members, eligible_count=len(eligible_members))

# 3. Add Member Route (Automatically assigns next queue position)
@app.route('/add', methods=['POST'])
def add_member():
    name = request.form.get('name')
    amount = float(request.form.get('amount'))
    
    # Calculate the next queue position
    total_members = Member.query.count()
    next_position = total_members + 1
    
    new_member = Member(name=name, monthly_amount=amount, queue_position=next_position)
    db.session.add(new_member)
    db.session.commit()
    
    return redirect(url_for('home'))

# 4. Route: Process Payout (Picks the next person in line)
@app.route('/draw', methods=['POST'])
def draw_winner():
    # Query the next person in line who hasn't received a payout yet
    next_in_line = Member.query.filter_by(has_received=False)\
                               .order_by(Member.queue_position.asc())\
                               .first()
    
    if next_in_line:
        next_in_line.has_received = True  # Mark as paid out
        db.session.commit()
        return f"<h1>🎉 Payout goes to Position #{next_in_line.queue_position}: {next_in_line.name}!</h1><a href='/'>Go Back</a>"
    else:
        return "<h1>All members in the queue have received their payout for this cycle!</h1><a href='/'>Go Back</a>"

# Route 5: Clear All Members from Queue
@app.route('/clear', methods=['POST'])
def clear_queue():
    # Delete all records in the Member table
    Member.query.delete()
    db.session.commit()
    return redirect(url_for('home'))    

# 5. Start the web server
if __name__ == '__main__':
    app.run(debug=True)

